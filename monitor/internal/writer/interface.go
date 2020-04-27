package writer

// Interface 监控写操作接口声明
type Interface interface {
	// WriteStr 写string消息
	// 参数
	// str消息(字符串类型)
	WriteStr(str string)
	// WriteByt 写byte消息
	// 参数
	// byts消息(字节类型)
	WriteByt(byts []byte)
	Start()
	Stop()
}
