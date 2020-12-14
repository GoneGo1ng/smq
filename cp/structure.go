package cp

import jsoniter "github.com/json-iterator/go"

// http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/errata01/os/mqtt-v3.1.1-errata01-os-complete.html
// 控制报文 Structure of an MQTT Control Packet
type ControlPacket struct {
	FixedHeader    *FixedHeader
	VariableHeader *VariableHeader
	Payload        *Payload
}

func (fh *ControlPacket) ToString() string {
	str, _ := jsoniter.MarshalToString(fh)
	return str
}

// 固定报头 Fixed header, present in all MQTT Control Packets
type FixedHeader struct {
	ControlPacketType uint8 // 报文类型
	Flags             *Flags
	RemainingLength   uint32 // 剩余长度
}

func (fh *FixedHeader) ToString() string {
	str, _ := jsoniter.MarshalToString(fh)
	return str
}

// 用于指定控制报文类型的标志位 Flags specific to each MQTT Control Packet type
type Flags struct {
	DUP    bool  // 控制报文的重复分发标志 Duplicate delivery of a PUBLISH Control Packet
	QoS    uint8 // PUBLISH报文的服务质量等级 PUBLISH Quality of Service
	Retain bool  // PUBLISH报文的保留标志 PUBLISH Retain flag
}

// 可变报头 Variable header, present in some MQTT Control Packets
type VariableHeader struct {
	TopicName        string // 主题名
	TopicLength      uint32 // 主题长度
	PacketIdentifier uint16 // 报文标识符
	Connect          *Connect
	ConnectAck       *ConnectAck
}

func (vh *VariableHeader) ToString() string {
	str, _ := jsoniter.MarshalToString(vh)
	return str
}

// 连接报文可变报
type Connect struct {
	ProtocolName  string // 协议名
	ProtocolLevel uint8  // 协议级别
	ConnectFlags  *ConnectFlags
	KeepAlive     uint16 // 保持连接
}

// 连接确认报文可变报头
type ConnectAck struct {
	AcknowledgeFlags *AcknowledgeFlags
	ReturnCode       uint8 // 连接返回码
}

// 连接确认标志
type AcknowledgeFlags struct {
	Reserved       uint8 // 保留位
	SessionPresent bool  // 当前会话
}

// 连接标志
type ConnectFlags struct {
	UsernameFlag bool  // 用户名标志
	PasswordFlag bool  // 密码标志
	WillRetain   bool  // 遗嘱保留
	WillQoS      uint8 // 遗嘱QoS
	WillFlag     bool  // 遗嘱标志
	CleanSession bool  // 清理会话
	Reserved     uint8 // 保留位
}

// 有效载荷 Payload, present in some MQTT Control Packets
type Payload struct {
	ConnectPayload   *ConnectPayload
	PublishPayload   *PublishPayload
	SubscribePayload *SubscribePayload
}

func (p *Payload) ToString() string {
	str, _ := jsoniter.MarshalToString(p)
	return str
}

// 连接载荷
type ConnectPayload struct {
	ClientIdentifier string // 客户端标识符
	WillTopic        string // 遗嘱主题
	WillMessage      string // 遗嘱消息
	Username         string // 用户名
	Password         string // 密码
}

// 发布载荷
type PublishPayload struct {
	Data []byte // 数据
}

// 订阅载荷
type SubscribePayload struct {
	TopicNameList    []string // 主题
	RequestedQoSList []uint8  // 服务质量要求
	ReturnCode       uint8    // 订阅返回码
}
