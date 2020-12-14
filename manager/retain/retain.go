package retain

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	jsoniter "github.com/json-iterator/go"
	"strings"
	"sync"
)

func init() {
	rr = &retainRoot{
		Nodes: make(map[string]*retainNode),
	}
}

var rr *retainRoot

type retainRoot struct {
	rmu   sync.RWMutex
	Nodes map[string]*retainNode
}

func GetRetainRoot() *retainRoot {
	return rr
}

func (rr *retainRoot) ShowRetainRoot() string {
	str, _ := jsoniter.MarshalToString(rr)
	return str
}

func (rr *retainRoot) PutRetainNode(packet *packets.PublishPacket) {
	rr.rmu.Lock()
	defer rr.rmu.Unlock()

	topic := packet.TopicName
	topicLevels := strings.Split(topic, "/")

	rn := rr.Nodes[topicLevels[0]]
	if rn == nil {
		rn = &retainNode{
			NextNodes: make(map[string]*retainNode),
		}
		rr.Nodes[topicLevels[0]] = rn
	}

	if len(topicLevels) == 1 {
		rn.Packet = packet
	} else {
		rn.putRetainNode(strings.Join(topicLevels[1:], "/"), packet)
	}
}

func (rr *retainRoot) DelRetainNode(topic string) {
	rr.rmu.Lock()
	defer rr.rmu.Unlock()

	topicLevels := strings.Split(topic, "/")

	rn := rr.Nodes[topicLevels[0]]
	if rn == nil {
		return
	}

	if len(topicLevels) == 1 {
		rn.Packet = nil
	} else {
		rn.delRetainNode(strings.Join(topicLevels[1:], "/"))
	}
	if rn.Packet == nil && len(rn.NextNodes) == 0 {
		delete(rr.Nodes, topicLevels[0])
	}
}

func (rr *retainRoot) MatchRetainNode(topic string) []packets.PublishPacket {
	rr.rmu.Lock()
	defer rr.rmu.Unlock()

	topicLevels := strings.Split(topic, "/")

	packets := &[]packets.PublishPacket{}

	rn := rr.Nodes[topicLevels[0]]

	if len(topicLevels) == 1 {
		if topic == "#" {
			for _, n := range rr.Nodes {
				if n.Packet != nil {
					*packets = append(*packets, *n.Packet)
				}
				n.matchRetainNode("#", packets)
			}
		} else if topic == "+" {
			for _, n := range rr.Nodes {
				if n.Packet != nil {
					*packets = append(*packets, *n.Packet)
				}
			}
		} else {
			if rn != nil && rn.Packet != nil {
				*packets = append(*packets, *rn.Packet)
			}
		}
	} else {
		if topicLevels[0] == "+" {
			for _, n := range rr.Nodes {
				n.matchRetainNode(strings.Join(topicLevels[1:], "/"), packets)
			}
		} else {
			if rn != nil {
				rn.matchRetainNode(strings.Join(topicLevels[1:], "/"), packets)
			}
		}
	}
	return *packets
}

type retainNode struct {
	Packet    *packets.PublishPacket
	NextNodes map[string]*retainNode
}

func (rn *retainNode) putRetainNode(topic string, packet *packets.PublishPacket) {
	topicLevels := strings.Split(topic, "/")
	if len(topicLevels) == 1 {
		if nn := rn.NextNodes[topic]; nn == nil {
			rn.NextNodes[topic] = &retainNode{
				Packet:    packet,
				NextNodes: make(map[string]*retainNode),
			}
		} else {
			nn.Packet = packet
		}
		return
	}
	nn := rn.NextNodes[topicLevels[0]]
	if nn == nil {
		nn = &retainNode{NextNodes: make(map[string]*retainNode)}
		rn.NextNodes[topicLevels[0]] = nn
	}
	nn.putRetainNode(strings.Join(topicLevels[1:], "/"), packet)
}

func (rn *retainNode) delRetainNode(topic string) {
	topicLevels := strings.Split(topic, "/")
	if len(topicLevels) == 1 {
		if nn := rn.NextNodes[topic]; nn != nil {
			nn.Packet = nil
			if nn.Packet == nil && len(nn.NextNodes) == 0 {
				delete(rn.NextNodes, topic)
			}
		}
		return
	}
	nn := rn.NextNodes[topicLevels[0]]
	if nn == nil {
		return
	}
	nn.delRetainNode(strings.Join(topicLevels[1:], "/"))
	if nn := rn.NextNodes[topic]; nn != nil {
		if nn.Packet == nil && len(nn.NextNodes) == 0 {
			delete(rn.NextNodes, topicLevels[0])
		}
	}
}

func (rn *retainNode) matchRetainNode(topic string, packets *[]packets.PublishPacket) {
	topicLevels := strings.Split(topic, "/")

	nn := rn.NextNodes[topicLevels[0]]

	if len(topicLevels) == 1 {
		if topic == "#" {
			for _, n := range rn.NextNodes {
				if n.Packet != nil {
					*packets = append(*packets, *n.Packet)
				}
				n.matchRetainNode("#", packets)
			}
		} else if topic == "+" {
			for _, n := range rn.NextNodes {
				if n.Packet != nil {
					*packets = append(*packets, *n.Packet)
				}
			}
		} else {
			if nn != nil && nn.Packet != nil {
				*packets = append(*packets, *nn.Packet)
			}
		}
		return
	}
	if topicLevels[0] == "+" {
		for _, n := range rn.NextNodes {
			n.matchRetainNode(strings.Join(topicLevels[1:], "/"), packets)
		}
	} else {
		if nn != nil {
			nn.matchRetainNode(strings.Join(topicLevels[1:], "/"), packets)
		}
	}
}
