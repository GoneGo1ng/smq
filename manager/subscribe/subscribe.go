package subscribe

import (
	jsoniter "github.com/json-iterator/go"
	"strings"
	"sync"
)

func init() {
	sr = &subscribeRoot{
		Nodes: make(map[string]*subscribeNode),
	}
}

var sr *subscribeRoot

type subscribeRoot struct {
	smu   sync.RWMutex
	Nodes map[string]*subscribeNode
}

func GetSubscribeRoot() *subscribeRoot {
	return sr
}

func (sr *subscribeRoot) ShowSubscribeRoot() string {
	str, _ := jsoniter.MarshalToString(sr)
	return str
}

func (sr *subscribeRoot) PutSubscribeNode(topic, clientID, groupName string, qos uint8, share bool) {
	sr.smu.Lock()
	defer sr.smu.Unlock()

	topicLevels := strings.Split(topic, "/")

	sn := sr.Nodes[topicLevels[0]]
	if sn == nil {
		sn = &subscribeNode{
			Subscriptions: make(map[string]*subscription),
			NextNodes:     make(map[string]*subscribeNode),
		}
		sr.Nodes[topicLevels[0]] = sn
	}

	if len(topicLevels) == 1 {
		sn.Subscriptions[clientID] = &subscription{
			ClientID:  clientID,
			Qos:       qos,
			Share:     share,
			GroupName: groupName,
		}
	} else {
		sn.putSubscribeNode(strings.Join(topicLevels[1:], "/"), clientID, groupName, qos, share)
	}
}

func (sr *subscribeRoot) DelSubscribeNode(topic, clientID string) {
	sr.smu.Lock()
	defer sr.smu.Unlock()

	topicLevels := strings.Split(topic, "/")

	sn := sr.Nodes[topicLevels[0]]
	if sn == nil {
		return
	}

	if len(topicLevels) == 1 {
		delete(sn.Subscriptions, clientID)
	} else {
		sn.delSubscribeNode(strings.Join(topicLevels[1:], "/"), clientID)
	}
	if len(sn.Subscriptions) == 0 && len(sn.NextNodes) == 0 {
		delete(sr.Nodes, topicLevels[0])
	}
}

func (sr *subscribeRoot) MatchSubscribeNode(topic string) []subscription {
	sr.smu.Lock()
	defer sr.smu.Unlock()

	topicLevels := strings.Split(topic, "/")

	subs := &[]subscription{}

	msn := sr.Nodes["#"]
	ssn := sr.Nodes["+"]
	sn := sr.Nodes[topicLevels[0]]

	if msn == nil && ssn == nil && sn == nil {
		return *subs
	}

	if msn != nil {
		for _, v := range msn.Subscriptions {
			*subs = append(*subs, *v)
		}
	}
	if len(topicLevels) == 1 {
		if ssn != nil {
			for _, v := range ssn.Subscriptions {
				*subs = append(*subs, *v)
			}
		}
		if sn != nil {
			for _, v := range sn.Subscriptions {
				*subs = append(*subs, *v)
			}
		}
	} else {
		/*if msn != nil {
			msn.matchSubscribeNode(strings.Join(topicLevels[1:], "/"), clientIDs)
		}*/
		if ssn != nil {
			ssn.matchSubscribeNode(strings.Join(topicLevels[1:], "/"), subs)
		}
		if sn != nil {
			sn.matchSubscribeNode(strings.Join(topicLevels[1:], "/"), subs)
		}
	}
	return *subs
}

func (sr *subscribeRoot) CleanSubscribeNode(clientID string) {
	sr.smu.Lock()
	defer sr.smu.Unlock()

	for topic, sn := range sr.Nodes {
		sn.cleanSubscribeNode(clientID)
		if len(sn.Subscriptions) == 0 && len(sn.NextNodes) == 0 {
			delete(sr.Nodes, topic)
		}
	}
}

type subscription struct {
	ClientID  string
	Qos       uint8
	Share     bool
	GroupName string
}

type subscribeNode struct {
	Subscriptions map[string]*subscription
	NextNodes     map[string]*subscribeNode
}

func (sn *subscribeNode) putSubscribeNode(topic, clientID, groupName string, qos uint8, share bool) {
	topicLevels := strings.Split(topic, "/")
	if len(topicLevels) == 1 {
		if nn := sn.NextNodes[topic]; nn == nil {
			sn.NextNodes[topic] = &subscribeNode{
				Subscriptions: map[string]*subscription{
					clientID: {
						ClientID:  clientID,
						Qos:       qos,
						Share:     share,
						GroupName: groupName,
					},
				},
				NextNodes: make(map[string]*subscribeNode),
			}
		} else {
			if nn.Subscriptions == nil {
				nn.Subscriptions = make(map[string]*subscription)
			}
			nn.Subscriptions[clientID] = &subscription{
				ClientID:  clientID,
				Qos:       qos,
				Share:     share,
				GroupName: groupName,
			}
		}
		return
	}
	nn := sn.NextNodes[topicLevels[0]]
	if nn == nil {
		nn = &subscribeNode{NextNodes: make(map[string]*subscribeNode)}
		sn.NextNodes[topicLevels[0]] = nn
	}
	nn.putSubscribeNode(strings.Join(topicLevels[1:], "/"), clientID, groupName, qos, share)
}

func (sn *subscribeNode) delSubscribeNode(topic, clientID string) {
	topicLevels := strings.Split(topic, "/")
	if len(topicLevels) == 1 {
		if nn := sn.NextNodes[topic]; nn != nil {
			delete(nn.Subscriptions, clientID)
			if len(nn.Subscriptions) == 0 && len(nn.NextNodes) == 0 {
				delete(sn.NextNodes, topic)
			}
		}
		return
	}
	nn := sn.NextNodes[topicLevels[0]]
	if nn == nil {
		return
	}
	nn.delSubscribeNode(strings.Join(topicLevels[1:], "/"), clientID)
	if nn := sn.NextNodes[topicLevels[0]]; nn != nil {
		if len(nn.Subscriptions) == 0 && len(nn.NextNodes) == 0 {
			delete(sn.NextNodes, topicLevels[0])
		}
	}
}

func (sn *subscribeNode) matchSubscribeNode(topic string, subs *[]subscription) {
	topicLevels := strings.Split(topic, "/")

	mnn := sn.NextNodes["#"]
	snn := sn.NextNodes["+"]
	nn := sn.NextNodes[topicLevels[0]]

	if mnn == nil && snn == nil && nn == nil {
		return
	}

	if mnn != nil {
		for _, v := range mnn.Subscriptions {
			*subs = append(*subs, *v)
		}
	}
	if len(topicLevels) == 1 {
		if snn != nil {
			for _, v := range snn.Subscriptions {
				*subs = append(*subs, *v)
			}
		}
		if nn != nil {
			for _, v := range nn.Subscriptions {
				*subs = append(*subs, *v)
			}
		}
		return
	}
	if snn != nil {
		snn.matchSubscribeNode(strings.Join(topicLevels[1:], "/"), subs)
	}
	if nn != nil {
		nn.matchSubscribeNode(strings.Join(topicLevels[1:], "/"), subs)
	}
}

func (sn *subscribeNode) cleanSubscribeNode(clientID string) {
	delete(sn.Subscriptions, clientID)
	for topic, nn := range sn.NextNodes {
		nn.cleanSubscribeNode(clientID)
		if len(nn.Subscriptions) == 0 && len(nn.NextNodes) == 0 {
			delete(sn.NextNodes, topic)
		}
	}
}
