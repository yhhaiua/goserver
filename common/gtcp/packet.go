package gtcp

// PacketBase 基础包结构
type PacketBase struct {
	Size    uint32
	Encrypt uint8
	keep    uint8
	Pcmd    BaseCmd
}

// BaseCmd 包头
type BaseCmd struct {
	Cmd    uint8
	SupCmd uint8
}

//Value 获取BasCmd的value值
func (data *BaseCmd) Value() uint16 {
	value := uint16(data.SupCmd) << 8
	value = value | uint16(data.Cmd)
	return value
}
