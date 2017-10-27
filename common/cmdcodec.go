package common

import "github.com/yhhaiua/goserver/common/goobjfmt"

//CmdCodec 包的解析接口
type CmdCodec interface {
	Encode(interface{}) ([]byte, error)

	Decode([]byte, interface{}) error

	Size(obj interface{}) int
}

//BinaryCodec 二进制包解析
type BinaryCodec struct {
}

//Encode 写入函数
func (codec *BinaryCodec) Encode(msgObj interface{}) ([]byte, error) {

	return goobjfmt.BinaryWrite(msgObj)

}

//Decode 读取函数
func (codec *BinaryCodec) Decode(data []byte, msgObj interface{}) error {

	return goobjfmt.BinaryRead(data, msgObj)
}

//Size 长度
func (codec *BinaryCodec) Size(msgObj interface{}) int {

	return goobjfmt.BinarySize(msgObj)
}
