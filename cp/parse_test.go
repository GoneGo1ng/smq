package cp

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseControlPacket(t *testing.T) {
	// CONNECT
	hexStream := "106900044d51545404c20078002d62393761303537663832316438656333386665363766363661653236356663617c323237386330306433393333000c32323738633030643339333300206462393235626562383034343931333963613433643730343466316137663363"
	buf, _ := hex.DecodeString(hexStream)
	cp := parseControlPacket(buf)
	fmt.Println(cp.ToString())

	// CONNACK
	hexStream = "20020000"
	buf, _ = hex.DecodeString(hexStream)
	cp = parseControlPacket(buf)
	fmt.Println(cp.ToString())

	// PUBLISH
	hexStream = "324000342f732f62393761303537663832316438656333386665363766363661653236356663612f3232373863303064333933332f652f6300137b2264223a7b7d7d"
	buf, _ = hex.DecodeString(hexStream)
	cp = parseControlPacket(buf)
	fmt.Println(cp.ToString())

	// PUBACK
	hexStream = "40020013"
	buf, _ = hex.DecodeString(hexStream)
	cp = parseControlPacket(buf)
	fmt.Println(cp.ToString())

	// SUBSCRIBE
	hexStream = "823e001400392f732f62393761303537663832316438656333386665363766363661653236356663612f3232373863303064333933332f722f2b2f72702f2b01"
	buf, _ = hex.DecodeString(hexStream)
	cp = parseControlPacket(buf)
	fmt.Println(cp.ToString())

	// SUBACK
	hexStream = "9003001401"
	buf, _ = hex.DecodeString(hexStream)
	cp = parseControlPacket(buf)
	fmt.Println(cp.ToString())

	// UNSUBSCRIBE
	hexStream = "a208000100042f666f6f"
	buf, _ = hex.DecodeString(hexStream)
	cp = parseControlPacket(buf)
	fmt.Println(cp.ToString())

	// UNSUBACK
	hexStream = "b0020001"
	buf, _ = hex.DecodeString(hexStream)
	cp = parseControlPacket(buf)
	fmt.Println(cp.ToString())

	// PINGREQ
	hexStream = "c000"
	buf, _ = hex.DecodeString(hexStream)
	cp = parseControlPacket(buf)
	fmt.Println(cp.ToString())

	// PINGRESP
	hexStream = "d000"
	buf, _ = hex.DecodeString(hexStream)
	cp = parseControlPacket(buf)
	fmt.Println(cp.ToString())

	// PINGRESP
	hexStream = "d000"
	buf, _ = hex.DecodeString(hexStream)
	cp = parseControlPacket(buf)
	fmt.Println(cp.ToString())

	// DISCONNECT
	hexStream = "e000"
	buf, _ = hex.DecodeString(hexStream)
	cp = parseControlPacket(buf)
	fmt.Println(cp.ToString())
}

func TestParseFixedHeader(t *testing.T) {
	// CONNECT
	hexStream := "106900044d51545404c20078002d62393761303537663832316438656333386665363766363661653236356663617c323237386330306433393333000c32323738633030643339333300206462393235626562383034343931333963613433643730343466316137663363"
	buf, _ := hex.DecodeString(hexStream)
	fh := parseFixedHeader(buf)
	fmt.Println(fh.ToString())
	assert.Equal(t, fh.ControlPacketType, CONNECT)
	assert.Equal(t, fh.RemainingLength, uint32(105))

	// CONNACK
	hexStream = "20020000"
	buf, _ = hex.DecodeString(hexStream)
	fh = parseFixedHeader(buf)
	fmt.Println(fh.ToString())
	assert.Equal(t, fh.ControlPacketType, CONNACK)
	assert.Equal(t, fh.RemainingLength, uint32(2))

	// PUBLISH
	hexStream = "324000342f732f62393761303537663832316438656333386665363766363661653236356663612f3232373863303064333933332f652f6300137b2264223a7b7d7d"
	buf, _ = hex.DecodeString(hexStream)
	fh = parseFixedHeader(buf)
	fmt.Println(fh.ToString())
	assert.Equal(t, fh.ControlPacketType, PUBLISH)
	assert.Equal(t, fh.Flags.DUP, false)
	assert.Equal(t, fh.Flags.QoS, uint8(1))
	assert.Equal(t, fh.Flags.Retain, false)
	assert.Equal(t, fh.RemainingLength, uint32(64))

	// PUBACK
	hexStream = "40020013"
	buf, _ = hex.DecodeString(hexStream)
	fh = parseFixedHeader(buf)
	fmt.Println(fh.ToString())
	assert.Equal(t, fh.ControlPacketType, PUBACK)
	assert.Equal(t, fh.RemainingLength, uint32(2))

	// SUBSCRIBE
	hexStream = "823e001400392f732f62393761303537663832316438656333386665363766363661653236356663612f3232373863303064333933332f722f2b2f72702f2b01"
	buf, _ = hex.DecodeString(hexStream)
	fh = parseFixedHeader(buf)
	fmt.Println(fh.ToString())
	assert.Equal(t, fh.ControlPacketType, SUBSCRIBE)
	assert.Equal(t, fh.Flags.DUP, false)
	assert.Equal(t, fh.Flags.QoS, uint8(1))
	assert.Equal(t, fh.Flags.Retain, false)
	assert.Equal(t, fh.RemainingLength, uint32(62))

	// SUBACK
	hexStream = "9003001401"
	buf, _ = hex.DecodeString(hexStream)
	fh = parseFixedHeader(buf)
	fmt.Println(fh.ToString())
	assert.Equal(t, fh.ControlPacketType, SUBACK)
	assert.Equal(t, fh.RemainingLength, uint32(3))

	// UNSUBSCRIBE
	hexStream = "a208000100042f666f6f"
	buf, _ = hex.DecodeString(hexStream)
	fh = parseFixedHeader(buf)
	fmt.Println(fh.ToString())
	assert.Equal(t, fh.ControlPacketType, UNSUBSCRIBE)
	assert.Equal(t, fh.RemainingLength, uint32(8))

	// UNSUBACK
	hexStream = "b0020001"
	buf, _ = hex.DecodeString(hexStream)
	fh = parseFixedHeader(buf)
	fmt.Println(fh.ToString())
	assert.Equal(t, fh.ControlPacketType, UNSUBACK)
	assert.Equal(t, fh.RemainingLength, uint32(2))

	// PINGREQ
	hexStream = "c000"
	buf, _ = hex.DecodeString(hexStream)
	fh = parseFixedHeader(buf)
	fmt.Println(fh.ToString())
	assert.Equal(t, fh.ControlPacketType, PINGREQ)
	assert.Equal(t, fh.RemainingLength, uint32(0))

	// PINGRESP
	hexStream = "d000"
	buf, _ = hex.DecodeString(hexStream)
	fh = parseFixedHeader(buf)
	fmt.Println(fh.ToString())
	assert.Equal(t, fh.ControlPacketType, PINGRESP)
	assert.Equal(t, fh.RemainingLength, uint32(0))

	// PINGRESP
	hexStream = "d000"
	buf, _ = hex.DecodeString(hexStream)
	fh = parseFixedHeader(buf)
	fmt.Println(fh.ToString())
	assert.Equal(t, fh.ControlPacketType, PINGRESP)
	assert.Equal(t, fh.RemainingLength, uint32(0))

	// DISCONNECT
	hexStream = "e000"
	buf, _ = hex.DecodeString(hexStream)
	fh = parseFixedHeader(buf)
	fmt.Println(fh.ToString())
	assert.Equal(t, fh.ControlPacketType, DISCONNECT)
	assert.Equal(t, fh.RemainingLength, uint32(0))
}

func TestParseVariableHeader(t *testing.T) {
	// CONNECT
	hexStream := "106900044d51545404c20078002d62393761303537663832316438656333386665363766363661653236356663617c323237386330306433393333000c32323738633030643339333300206462393235626562383034343931333963613433643730343466316137663363"
	buf, _ := hex.DecodeString(hexStream)
	fh := parseFixedHeader(buf)
	fmt.Println(fh.ToString())
	vh, i := parseVariableHeader(fh, buf)
	fmt.Println(i, vh.ToString())
	assert.Equal(t, i, 12)
	assert.Equal(t, vh.Connect.ProtocolName, "MQTT")
	assert.Equal(t, vh.Connect.ProtocolLevel, uint8(4))
	assert.Equal(t, vh.Connect.ConnectFlags.UsernameFlag, true)
	assert.Equal(t, vh.Connect.ConnectFlags.PasswordFlag, true)
	assert.Equal(t, vh.Connect.ConnectFlags.WillRetain, false)
	assert.Equal(t, vh.Connect.ConnectFlags.WillQoS, uint8(0))
	assert.Equal(t, vh.Connect.ConnectFlags.CleanSession, true)
	assert.Equal(t, vh.Connect.ConnectFlags.Reserved, uint8(0))
	assert.Equal(t, vh.Connect.KeepAlive, uint16(120))

	// CONNACK
	hexStream = "20020000"
	buf, _ = hex.DecodeString(hexStream)
	fh = parseFixedHeader(buf)
	fmt.Println(fh.ToString())
	vh, i = parseVariableHeader(fh, buf)
	fmt.Println(i, vh.ToString())
	assert.Equal(t, i, 4)
	assert.Equal(t, vh.ConnectAck.AcknowledgeFlags.Reserved, uint8(0))
	assert.Equal(t, vh.ConnectAck.AcknowledgeFlags.SessionPresent, false)
	assert.Equal(t, vh.ConnectAck.ReturnCode, uint8(0))

	// PUBLISH
	hexStream = "324000342f732f62393761303537663832316438656333386665363766363661653236356663612f3232373863303064333933332f652f6300137b2264223a7b7d7d"
	buf, _ = hex.DecodeString(hexStream)
	fh = parseFixedHeader(buf)
	fmt.Println(fh.ToString())
	vh, i = parseVariableHeader(fh, buf)
	fmt.Println(i, vh.ToString())
	assert.Equal(t, i, 58)
	assert.Equal(t, vh.TopicName, "/s/b97a057f821d8ec38fe67f66ae265fca/2278c00d3933/e/c")
	assert.Equal(t, vh.PacketIdentifier, uint16(19))

	// PUBACK
	hexStream = "40020013"
	buf, _ = hex.DecodeString(hexStream)
	fh = parseFixedHeader(buf)
	fmt.Println(fh.ToString())
	vh, i = parseVariableHeader(fh, buf)
	fmt.Println(i, vh.ToString())
	assert.Equal(t, vh.PacketIdentifier, uint16(19))

	// SUBSCRIBE
	hexStream = "823e001400392f732f62393761303537663832316438656333386665363766363661653236356663612f3232373863303064333933332f722f2b2f72702f2b01"
	buf, _ = hex.DecodeString(hexStream)
	fh = parseFixedHeader(buf)
	fmt.Println(fh.ToString())
	vh, i = parseVariableHeader(fh, buf)
	fmt.Println(i, vh.ToString())
	assert.Equal(t, vh.PacketIdentifier, uint16(20))

	// SUBACK
	hexStream = "9003001401"
	buf, _ = hex.DecodeString(hexStream)
	fh = parseFixedHeader(buf)
	fmt.Println(fh.ToString())
	vh, i = parseVariableHeader(fh, buf)
	fmt.Println(i, vh.ToString())
	assert.Equal(t, vh.PacketIdentifier, uint16(20))

	// UNSUBSCRIBE
	hexStream = "a208000100042f666f6f"
	buf, _ = hex.DecodeString(hexStream)
	fh = parseFixedHeader(buf)
	fmt.Println(fh.ToString())
	vh, i = parseVariableHeader(fh, buf)
	fmt.Println(i, vh.ToString())
	assert.Equal(t, vh.PacketIdentifier, uint16(1))

	// UNSUBACK
	hexStream = "b0020001"
	buf, _ = hex.DecodeString(hexStream)
	fh = parseFixedHeader(buf)
	fmt.Println(fh.ToString())
	vh, i = parseVariableHeader(fh, buf)
	fmt.Println(i, vh.ToString())
	assert.Equal(t, vh.PacketIdentifier, uint16(1))
}

func TestParsePayload(t *testing.T) {
	// CONNECT
	hexStream := "106900044d51545404c20078002d62393761303537663832316438656333386665363766363661653236356663617c323237386330306433393333000c32323738633030643339333300206462393235626562383034343931333963613433643730343466316137663363"
	buf, _ := hex.DecodeString(hexStream)
	fh := parseFixedHeader(buf)
	fmt.Println(fh.ToString())
	vh, i := parseVariableHeader(fh, buf)
	fmt.Println(i, vh.ToString())
	p := parsePayload(fh, vh, buf, i)
	fmt.Println(p.ToString())
	assert.Equal(t, p.ConnectPayload.ClientIdentifier, "b97a057f821d8ec38fe67f66ae265fca|2278c00d3933")
	assert.Equal(t, p.ConnectPayload.WillTopic, "")
	assert.Equal(t, p.ConnectPayload.WillMessage, "")
	assert.Equal(t, p.ConnectPayload.Username, "2278c00d3933")
	assert.Equal(t, p.ConnectPayload.Password, "db925beb80449139ca43d7044f1a7f3c")

	// PUBLISH
	hexStream = "324000342f732f62393761303537663832316438656333386665363766363661653236356663612f3232373863303064333933332f652f6300137b2264223a7b7d7d"
	buf, _ = hex.DecodeString(hexStream)
	fh = parseFixedHeader(buf)
	fmt.Println(fh.ToString())
	vh, i = parseVariableHeader(fh, buf)
	fmt.Println(i, vh.ToString())
	p = parsePayload(fh, vh, buf, i)
	fmt.Println(p.ToString())
	assert.Equal(t, string(p.PublishPayload.Data), `{"d":{}}`)

	// SUBSCRIBE
	hexStream = "823e001400392f732f62393761303537663832316438656333386665363766363661653236356663612f3232373863303064333933332f722f2b2f72702f2b01"
	buf, _ = hex.DecodeString(hexStream)
	fh = parseFixedHeader(buf)
	fmt.Println(fh.ToString())
	vh, i = parseVariableHeader(fh, buf)
	fmt.Println(i, vh.ToString())
	p = parsePayload(fh, vh, buf, i)
	fmt.Println(p.ToString())
	assert.Equal(t, p.SubscribePayload.TopicNameList[0], "/s/b97a057f821d8ec38fe67f66ae265fca/2278c00d3933/r/+/rp/+")
	assert.Equal(t, p.SubscribePayload.RequestedQoSList[0], uint8(1))

	// SUBACK
	hexStream = "9003001401"
	buf, _ = hex.DecodeString(hexStream)
	fh = parseFixedHeader(buf)
	fmt.Println(fh.ToString())
	vh, i = parseVariableHeader(fh, buf)
	fmt.Println(i, vh.ToString())
	p = parsePayload(fh, vh, buf, i)
	fmt.Println(p.ToString())
	assert.Equal(t, p.SubscribePayload.ReturnCode, uint8(1))

	// UNSUBSCRIBE
	hexStream = "a208000100042f666f6f"
	buf, _ = hex.DecodeString(hexStream)
	fh = parseFixedHeader(buf)
	fmt.Println(fh.ToString())
	vh, i = parseVariableHeader(fh, buf)
	fmt.Println(i, vh.ToString())
	p = parsePayload(fh, vh, buf, i)
	fmt.Println(p.ToString())
	assert.Equal(t, p.SubscribePayload.TopicNameList[0], "/foo")
}
