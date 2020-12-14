package cp

const (
	CONNECT = uint8(iota + 1)
	CONNACK
	PUBLISH
	PUBACK
	PUBREC
	PUBREL
	PUBCOMP
	SUBSCRIBE
	SUBACK
	UNSUBSCRIBE
	UNSUBACK
	PINGREQ
	PINGRESP
	DISCONNECT
)

func parseControlPacket(buf []byte) *ControlPacket {
	fh := parseFixedHeader(buf)
	vh, i := parseVariableHeader(fh, buf)
	p := parsePayload(fh, vh, buf, i)
	cp := &ControlPacket{}
	cp.FixedHeader = fh
	cp.VariableHeader = vh
	cp.Payload = p
	return cp
}

func parseFixedHeader(buf []byte) *FixedHeader {
	fh := &FixedHeader{}
	fh.ControlPacketType = buf[0] & 0xF0 >> 4
	if fh.ControlPacketType == PUBLISH || fh.ControlPacketType == PUBREL ||
		fh.ControlPacketType == SUBSCRIBE || fh.ControlPacketType == UNSUBSCRIBE {
		f := &Flags{}
		f.DUP = buf[0]&0x08 > 0
		f.QoS = buf[0] & 0x06 >> 1
		f.Retain = buf[0]&0x01 > 0
		fh.Flags = f
	}
	fh.RemainingLength = uint32(buf[1])
	return fh
}

func parseVariableHeader(fh *FixedHeader, buf []byte) (*VariableHeader, int) {
	i := 2
	vh := &VariableHeader{}
	switch fh.ControlPacketType {
	case CONNECT:
		c := &Connect{}
		c.ProtocolName = getString(buf, &i)
		c.ProtocolLevel = getUint8(buf, &i)
		cfb := getByte(buf, &i)
		cf := &ConnectFlags{}
		cf.UsernameFlag = cfb&0x80 > 0
		cf.PasswordFlag = cfb&0x40 > 0
		cf.WillRetain = cfb&0x20 > 0
		cf.WillQoS = cfb & 0x18 >> 3
		cf.CleanSession = cfb&0x02 > 0
		cf.Reserved = 0
		c.ConnectFlags = cf
		c.KeepAlive = getUint16(buf, &i)
		vh.Connect = c
	case CONNACK:
		ca := &ConnectAck{}
		af := &AcknowledgeFlags{}
		afb := getByte(buf, &i)
		af.Reserved = 0
		af.SessionPresent = afb&0x01 > 0
		ca.AcknowledgeFlags = af
		ca.ReturnCode = getUint8(buf, &i)
		vh.ConnectAck = ca
	case PUBLISH:
		tn, tl := getStringAndLength(buf, &i)
		vh.TopicName = tn
		vh.TopicLength = uint32(tl)
		if qos := fh.Flags.QoS; qos == 1 || qos == 2 {
			vh.PacketIdentifier = getUint16(buf, &i)
		}
	case PUBACK, PUBREC, PUBREL, PUBCOMP, SUBSCRIBE, SUBACK, UNSUBSCRIBE, UNSUBACK:
		vh.PacketIdentifier = getUint16(buf, &i)
	}
	return vh, i
}

func parsePayload(fh *FixedHeader, vh *VariableHeader, buf []byte, i int) *Payload {
	p := &Payload{}
	switch fh.ControlPacketType {
	case CONNECT:
		cp := &ConnectPayload{}
		cp.ClientIdentifier = getString(buf, &i)
		if vh.Connect.ConnectFlags.WillFlag {
			cp.WillTopic = getString(buf, &i)
			cp.WillMessage = getString(buf, &i)
		}
		cp.Username = getString(buf, &i)
		cp.Password = getString(buf, &i)
		p.ConnectPayload = cp
	case PUBLISH:
		pp := &PublishPayload{}
		pp.Data = getBytes(buf, &i, fh.RemainingLength-vh.TopicLength-4)
		p.PublishPayload = pp
	case SUBSCRIBE:
		sp := &SubscribePayload{}
		for ; i < len(buf); {
			sp.TopicNameList = append(sp.TopicNameList, getString(buf, &i))
			sp.RequestedQoSList = append(sp.RequestedQoSList, getUint8(buf, &i))
		}
		p.SubscribePayload = sp
	case SUBACK:
		sp := &SubscribePayload{}
		sp.ReturnCode = getUint8(buf, &i)
		p.SubscribePayload = sp
	case UNSUBSCRIBE:
		sp := &SubscribePayload{}
		for ; i < len(buf); {
			sp.TopicNameList = append(sp.TopicNameList, getString(buf, &i))
		}
		p.SubscribePayload = sp
	}
	return p
}

func getByte(b []byte, i *int) byte {
	*i += 1
	return b[*i-1]
}

func getUint8(b []byte, i *int) uint8 {
	*i += 1
	return b[*i-1]
}

func getUint16(b []byte, i *int) uint16 {
	*i += 2
	return uint16(b[*i-2])<<8 + uint16(b[*i-1])
}

func getString(b []byte, i *int) string {
	l := int(getUint16(b, i))
	*i += l
	return string(b[*i-l : *i])
}

func getStringAndLength(b []byte, i *int) (string, int) {
	l := int(getUint16(b, i))
	*i += l
	return string(b[*i-l : *i]), l
}

func getBytes(b []byte, i *int, l uint32) []byte {
	return b[*i : *i+int(l)]
}
