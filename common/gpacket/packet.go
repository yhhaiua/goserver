package gpacket

// PacketBase 基础包结构
type PacketBase struct {
	Size    uint32
	Encrypt uint8
	Keep    uint8
	Pcmd    BaseCmd
}

// BaseCmd 包头
type BaseCmd struct {
	Cmd    uint8
	SupCmd uint8
}

// Packet 一个包结构
type Packet struct {
	Size    uint32
	Encrypt uint8
	Keep    uint8
}

//Value 获取BasCmd的value值
func (data *BaseCmd) Value() uint16 {
	value := uint16(data.SupCmd) << 8
	value = value | uint16(data.Cmd)
	return value
}

//GetValue 获取value值
func GetValue(Cmd, SupCmd uint8) uint16 {
	value := uint16(SupCmd) << 8
	value = value | uint16(Cmd)
	return value
}
