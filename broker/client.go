package broker

import (
	"bytes"
	"context"
	"errors"
	"github.com/Allenxuxu/gev/connection"
	"github.com/eclipse/paho.mqtt.golang/packets"
	"go.uber.org/zap"
	"github.com/GoneGo1ng/smq/auth"
	"github.com/GoneGo1ng/smq/manager/retain"
	"github.com/GoneGo1ng/smq/manager/subscribe"
	"regexp"
	"sync"
	"sync/atomic"
	"time"
)

const (
	Connected = iota
	Disconnected
)

const (
	Qos0 = iota
	Qos1
	Qos2
)

const (
	GroupTopicRegexp = `^\$share/([0-9a-zA-Z_-]+)/(.*)$`
)

var (
	groupCompile = regexp.MustCompile(GroupTopicRegexp)
)

type client struct {
	mu             sync.RWMutex
	conn           *connection.Connection
	info           info
	status         int
	ctx            context.Context
	cancelFunc     context.CancelFunc
	noAckMsgID     uint32
	publishMsgs    sync.Map
	retainMsgs     sync.Map
	offlineMsgs    sync.Map
	keepaliveTimer *time.Timer
}

type info struct {
	clientID     string
	username     string
	password     []byte
	keepalive    uint16
	timeout      time.Duration
	willMsg      *packets.PublishPacket
	cleanSession bool
	localIP      string
	remoteIP     string
}

func (c *client) init() {
	c.status = Connected
	c.info.localIP = c.conn.PeerAddr()
	c.ctx, c.cancelFunc = context.WithCancel(context.Background())
	go c.publishNoAck()
	go c.checkAlive()
}

func (c *client) processPacket(packet packets.ControlPacket) int {
	if packet == nil {
		return 0
	}
	if c.info.timeout != 0 {
		c.keepaliveTimer.Reset(c.info.timeout)
	}
	switch packet.(type) {
	case *packets.ConnectPacket:
		return packet.(*packets.ConnectPacket).RemainingLength
	case *packets.ConnackPacket:
		return packet.(*packets.ConnackPacket).RemainingLength
	case *packets.PublishPacket:
		return c.processPublish(packet.(*packets.PublishPacket), false)
	case *packets.PubackPacket:
		return c.processPuback(packet.(*packets.PubackPacket))
	case *packets.PubrecPacket:
		return packet.(*packets.PubrecPacket).RemainingLength
	case *packets.PubrelPacket:
		return packet.(*packets.PubrelPacket).RemainingLength
	case *packets.PubcompPacket:
		return packet.(*packets.PubcompPacket).RemainingLength
	case *packets.SubscribePacket:
		return c.processSubscribe(packet.(*packets.SubscribePacket))
	case *packets.SubackPacket:
		return packet.(*packets.SubackPacket).RemainingLength
	case *packets.UnsubscribePacket:
		return c.processUnSubscribe(packet.(*packets.UnsubscribePacket))
	case *packets.UnsubackPacket:
		return packet.(*packets.UnsubackPacket).RemainingLength
	case *packets.PingreqPacket:
		return c.processPing(packet.(*packets.PingreqPacket))
	case *packets.PingrespPacket:
		return packet.(*packets.PingrespPacket).RemainingLength
	case *packets.DisconnectPacket:
		c.status = Disconnected
		c.conn.Close()
		return packet.(*packets.DisconnectPacket).RemainingLength
	default:
		c.conn.Close()
		zap.L().Warn("未知请求", zap.String("ClientID", c.info.clientID), zap.String("Packet", packet.String()))
		return 0
	}
}

func (c *client) processConnect(packet *packets.ConnectPacket) int {
	ca := packets.NewControlPacket(packets.Connack).(*packets.ConnackPacket)
	ca.ReturnCode = packet.Validate()

	if err := c.send(ca); err != nil {
		zap.L().Error("响应失败", zap.String("ClientID", packet.ClientIdentifier), zap.Error(err))
		return packet.RemainingLength
	}
	if ca.ReturnCode != packets.Accepted {
		zap.L().Error("未接受的连接", zap.String("ClientID", packet.ClientIdentifier),
			zap.Any("ReturnCode", ca.ReturnCode))
		c.conn.Close()
		return packet.RemainingLength
	}

	if _, ok := c.conn.Context().(*client); ok {
		zap.L().Warn("重复发送连接请求", zap.String("ClientID", packet.ClientIdentifier))
		c.conn.Close()
		return packet.RemainingLength
	}

	will := packets.NewControlPacket(packets.Publish).(*packets.PublishPacket)
	if packet.WillFlag {
		will.Qos = packet.WillQos
		will.TopicName = packet.WillTopic
		will.Retain = packet.WillRetain
		will.Payload = packet.WillMessage
		will.Dup = packet.Dup
	} else {
		will = nil
	}
	c.info = info{
		clientID:     packet.ClientIdentifier,
		username:     packet.Username,
		password:     packet.Password,
		keepalive:    packet.Keepalive,
		timeout:      time.Duration(packet.Keepalive+packet.Keepalive/2) * time.Second,
		willMsg:      will,
		cleanSession: packet.CleanSession,
	}
	if c.info.timeout != 0 {
		c.keepaliveTimer = time.NewTimer(c.info.timeout)
	}
	c.init()
	return packet.RemainingLength
}

func (c *client) processPublish(packet *packets.PublishPacket, will bool) int {
	if !auth.CheckTopicAuth(auth.PUB, c.info.clientID, c.info.username, packet.TopicName) {
		zap.L().Error("权限校验失败", zap.String("ClientID", c.info.clientID), zap.String("Topic", packet.TopicName))
		c.conn.Close()
		return packet.RemainingLength
	}

	switch packet.Qos {
	case Qos0:
		c.publish(packet)
		return packet.RemainingLength
	case Qos1:
		if !will {
			pap := packets.NewControlPacket(packets.Puback).(*packets.PubackPacket)
			pap.MessageID = packet.MessageID
			if err := c.send(pap); err != nil {
				zap.L().Error("响应失败", zap.String("ClientID", c.info.clientID), zap.Error(err))
				return packet.RemainingLength
			}
		}
		c.publish(packet)
		return packet.RemainingLength
	case Qos2:
		// TODO Qos2
		zap.L().Warn("暂不支持QOS2", zap.String("ClientID", c.info.clientID))
		return packet.RemainingLength
	default:
		zap.L().Error("未知的QOS", zap.String("ClientID", c.info.clientID))
		c.conn.Close()
		return packet.RemainingLength
	}
}

func (c *client) processPuback(packet *packets.PubackPacket) int {
	c.publishMsgs.Delete(packet.MessageID)
	c.offlineMsgs.Delete(packet.MessageID)
	// TODO 是否启用Retain参数
	c.retainMsgs.Delete(packet.MessageID)
	return packet.RemainingLength
}

func (c *client) publish(packet *packets.PublishPacket) {
	gm := map[string]bool{}
	subs := subscribe.GetSubscribeRoot().MatchSubscribeNode(packet.TopicName)
	broker.clients.Range(func(key, value interface{}) bool {
		for _, sub := range subs {
			if sub.ClientID == key.(string) {
				// 判断分组
				// zap.L().Debug("sub", zap.Any("sub", sub), zap.Any("gm", gm))
				if sub.Share {
					if !gm[sub.GroupName] {
						gm[sub.GroupName] = true
					} else {
						continue
					}
				}
				// 组织报文
				if packet.Qos > sub.Qos {
					packet.Qos = sub.Qos
				}
				client := value.(*client)
				if client.conn == nil {
					if packet.Qos == 1 {
						atomic.AddUint32(&client.noAckMsgID, 1)
						packet.MessageID = uint16(client.noAckMsgID)
						client.offlineMsgs.Store(packet.MessageID, *packet)
					}
					continue
				}
				if packet.Qos == 1 {
					atomic.AddUint32(&client.noAckMsgID, 1)
					packet.MessageID = uint16(client.noAckMsgID)
					client.publishMsgs.Store(packet.MessageID, *packet)
				} else if packet.Qos == 2 {
					// TODO Qos2 Puback
				}
				if err := client.send(packet); err != nil {
					zap.L().Error("响应失败", zap.String("clientID", c.info.clientID), zap.Error(err))
				}
			}
		}
		return true
	})

	// TODO 是否启用Retain参数
	// Retain
	if packet.Retain {
		if len(packet.Payload) == 0 {
			retain.GetRetainRoot().DelRetainNode(packet.TopicName)
		} else {
			retain.GetRetainRoot().PutRetainNode(packet)
		}
	}
}

func (c *client) processSubscribe(packet *packets.SubscribePacket) int {
	sap := packets.NewControlPacket(packets.Suback).(*packets.SubackPacket)
	sap.MessageID = packet.MessageID

	for i, topic := range packet.Topics {
		if !auth.CheckTopicAuth(auth.PUB, c.info.clientID, c.info.username, topic) {
			zap.L().Error("权限校验失败", zap.String("ClientID", c.info.clientID), zap.String("topic", topic))
			c.conn.Close()
			return packet.RemainingLength
		}
		subs := groupCompile.FindStringSubmatch(topic)
		if len(subs) == 0 {
			subscribe.GetSubscribeRoot().
				PutSubscribeNode(topic, c.info.clientID, "", packet.Qoss[i], false)
		} else if len(subs) == 3 {
			// $share
			subscribe.GetSubscribeRoot().
				PutSubscribeNode(subs[2], c.info.clientID, subs[1], packet.Qoss[i], true)
		} else {
			zap.L().Error("订阅主题错误", zap.String("ClientID", c.info.clientID), zap.String("topic", topic))
			c.conn.Close()
			return packet.RemainingLength
		}
		sap.ReturnCodes = append(sap.ReturnCodes, packet.Qoss[i])
	}

	if err := c.send(sap); err != nil {
		zap.L().Error("响应失败", zap.String("ClientID", c.info.clientID), zap.Error(err))
	}

	// TODO 是否启用Retain参数
	// Retain
	for i, topic := range packet.Topics {
		pps := retain.GetRetainRoot().MatchRetainNode(topic)
		for _, pp := range pps {
			pp.Qos = packet.Qoss[i]
			atomic.AddUint32(&c.noAckMsgID, 1)
			pp.MessageID = uint16(c.noAckMsgID)
			c.retainMsgs.Store(pp.MessageID, pp)
			if err := c.send(&pp); err != nil {
				zap.L().Error("响应失败", zap.String("ClientID", c.info.clientID), zap.Error(err))
			}
		}
	}
	return packet.RemainingLength
}

func (c *client) processUnSubscribe(packet *packets.UnsubscribePacket) int {
	uap := packets.NewControlPacket(packets.Unsuback).(*packets.UnsubackPacket)
	uap.MessageID = packet.MessageID
	if err := c.send(uap); err != nil {
		zap.L().Error("响应失败", zap.String("ClientID", c.info.clientID), zap.Error(err))
		return packet.RemainingLength
	}

	for _, topic := range packet.Topics {
		subscribe.GetSubscribeRoot().DelSubscribeNode(topic, c.info.clientID)
	}
	return packet.RemainingLength
}

func (c *client) processPing(packet *packets.PingreqPacket) int {
	pap := packets.NewControlPacket(packets.Pingresp).(*packets.PingrespPacket)
	if err := c.send(pap); err != nil {
		zap.L().Error("响应失败", zap.String("ClientID", c.info.clientID), zap.Error(err))
	}
	return packet.RemainingLength
}

func (c *client) send(packet packets.ControlPacket) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn == nil {
		return errors.New("conn is nil")
	}
	buf := new(bytes.Buffer)
	if err := packet.Write(buf); err != nil {
		return err
	}
	if err := c.conn.Send(buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (c *client) publishNoAck() {
	// TODO 间隔时间参数
	for range time.Tick(20 * time.Second) {
		select {
		case <-c.ctx.Done():
			return
		default:
			/*if c.conn == nil || !c.conn.Connected() || c.status == Disconnected {
				return
			}*/
			c.publishMsgs.Range(func(key, value interface{}) bool {
				r := value.(packets.PublishPacket)
				if err := c.send(&r); err != nil {
					zap.L().Error("响应失败", zap.String("ClientID", c.info.clientID), zap.Error(err))
				}
				return true
			})
			c.offlineMsgs.Range(func(key, value interface{}) bool {
				r := value.(packets.PublishPacket)
				if err := c.send(&r); err != nil {
					zap.L().Error("响应失败", zap.String("ClientID", c.info.clientID), zap.Error(err))
				}
				return true
			})
			// TODO 是否启用Retain参数
			c.retainMsgs.Range(func(key, value interface{}) bool {
				r := value.(packets.PublishPacket)
				if err := c.send(&r); err != nil {
					zap.L().Error("响应失败", zap.String("ClientID", c.info.clientID), zap.Error(err))
				}
				return true
			})
		}
	}
}

func (c *client) checkAlive() {
	if c.info.keepalive == 0 {
		return
	}
	select {
	case <-c.ctx.Done():
		c.keepaliveTimer.Stop()
		return
	case <-c.keepaliveTimer.C:
		zap.L().Warn("连接不活跃，关闭连接", zap.String("ClientID", c.info.clientID))
		c.conn.Close()
		return
	}
}
