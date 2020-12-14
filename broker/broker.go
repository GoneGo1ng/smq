package broker

import (
	"bytes"
	"github.com/Allenxuxu/gev/connection"
	"github.com/aristanetworks/goarista/monotime"
	"github.com/eclipse/paho.mqtt.golang/packets"
	"go.uber.org/zap"
	"hzsun.com/easyconnect/smq/manager/subscribe"
	"sync"
)

type Broker struct {
	clients sync.Map
}

var broker *Broker

func LoadBroker() (*Broker, error) {
	if broker == nil {
		broker = new(Broker)
		return broker, nil
	}
	return broker, nil
}

func (b *Broker) OnConnect(c *connection.Connection) {
	zap.L().Info("OnConnect", zap.String("PeerAddr", c.PeerAddr()))
}

func (b *Broker) OnMessage(c *connection.Connection, ctx interface{}, data []byte) (out []byte) {
	messageId := int(monotime.Now())
	zap.L().Info("OnMessage", zap.Int("MessageId", messageId), zap.ByteString("Data", data))
	for len(data) > 0 {
		packet, err := packets.ReadPacket(bytes.NewBuffer(data))
		if err != nil {
			zap.L().Error("解析报文失败", zap.Int("MessageId", messageId), zap.Error(err))
			c.Close()
			return
		}
		if packet == nil {
			zap.L().Error("空报文", zap.Int("MessageId", messageId))
			c.Close()
			return
		}
		zap.L().Debug("解析报文成功", zap.Int("MessageId", messageId), zap.String("Packet", packet.String()))

		if cp, ok := packet.(*packets.ConnectPacket); ok {
			nc := &client{conn: c}
			if rl := nc.processConnect(cp); rl != 0 {
				data = data[rl+2:]
			} else {
				return
			}
			if v, ok := b.clients.Load(cp.ClientIdentifier); ok {
				if !cp.CleanSession {
					if oc, ok := v.(*client); ok {
						nc.noAckMsgID = oc.noAckMsgID
						oc.offlineMsgs.Range(func(key, value interface{}) bool {
							pp := value.(packets.PublishPacket)
							nc.offlineMsgs.Store(key, value)
							if err := nc.send(&pp); err != nil {
								zap.L().Error("响应失败", zap.String("ClientID", nc.info.clientID), zap.Error(err))
							}
							return true
						})
					}
				} else {
					subscribe.GetSubscribeRoot().CleanSubscribeNode(nc.info.clientID)
				}
			}
			c.SetContext(nc)
			b.clients.Store(cp.ClientIdentifier, nc)
		} else {
			if client, ok := c.Context().(*client); ok {
				if rl := client.processPacket(packet); rl != 0 {
					data = data[rl+2:]
				} else {
					return
				}
			} else {
				zap.L().Warn("连接不存在", zap.Int("MessageId", messageId),
					zap.String("Packet", packet.String()))
				c.Close()
			}
		}
	}
	return
}

func (b *Broker) OnClose(c *connection.Connection) {
	if client, ok := c.Context().(*client); ok {
		// will
		if client.info.willMsg != nil && client.status == Connected {
			zap.L().Debug("发送遗嘱", zap.String("ClientID", client.info.clientID),
				zap.String("Will", client.info.willMsg.String()))
			client.processPublish(client.info.willMsg, true)
		}
		if client.info.cleanSession {
			b.clients.Delete(client.info.clientID)
			// clean subscribe topic
			subscribe.GetSubscribeRoot().CleanSubscribeNode(client.info.clientID)
		} else {
			client.status = Disconnected
			client.conn = nil
		}
		client.cancelFunc()
		zap.L().Info("离线", zap.String("ClientID", client.info.clientID))
	}
}
