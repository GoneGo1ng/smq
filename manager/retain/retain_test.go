package retain

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang/packets"
	"testing"
)

func TestPutRetainNode(t *testing.T) {
	rr := GetRetainRoot()
	p := packets.NewControlPacket(packets.Publish).(*packets.PublishPacket)
	p.TopicName = "foo"
	rr.PutRetainNode(p)
	fmt.Println(rr.ShowRetainRoot())
}

func TestDelRetainNode(t *testing.T) {
	rr := GetRetainRoot()
	p := packets.NewControlPacket(packets.Publish).(*packets.PublishPacket)
	p.TopicName = "foo"
	rr.PutRetainNode(p)
	p1 := packets.NewControlPacket(packets.Publish).(*packets.PublishPacket)
	p1.TopicName = "foo1"
	rr.PutRetainNode(p1)
	fmt.Println(rr.ShowRetainRoot())

	rr.DelRetainNode("foo")
	fmt.Println(rr.ShowRetainRoot())
}

func TestMatchRetainNode(t *testing.T) {
	rr := GetRetainRoot()
	p := packets.NewControlPacket(packets.Publish).(*packets.PublishPacket)
	p.TopicName = "/a"
	rr.PutRetainNode(p)
	p1 := packets.NewControlPacket(packets.Publish).(*packets.PublishPacket)
	p1.TopicName = "/a/b"
	rr.PutRetainNode(p1)
	p2 := packets.NewControlPacket(packets.Publish).(*packets.PublishPacket)
	p2.TopicName = "/a/b/c"
	rr.PutRetainNode(p2)
	fmt.Println(rr.ShowRetainRoot())

	ps := rr.MatchRetainNode("#")
	fmt.Println(ps)
}