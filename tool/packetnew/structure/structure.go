package structure

import (
	"encoding/xml"
	"log"
)

//StField field
type StField struct {
	Name    string `xml:"name,attr"`
	Type    string `xml:"type,attr"`
	RefType string `xml:"refType,attr"`
	Des     string `xml:"des,attr"`
}

//OnePacketGet 结构
type OnePacketGet struct {
	Id    string    `xml:"id,attr"`
	Name   string    `xml:"name,attr"`
	Type string    `xml:"type,attr"`
	Des    string    `xml:"des,attr"`
	Field  []StField `xml:"field"`
}

//PacketGet 结构
type PacketGet struct {
	Packet []OnePacketGet `xml:"packet"`
}

//OneRef 结构
type OneRef struct {
	Name  string    `xml:"name,attr"`
	Des   string    `xml:"des,attr"`
	Field []StField `xml:"field"`
}

//Ref 结构
type Ref struct {
	MyRef []OneRef `xml:"ref"`
}

var MyPacketGet PacketGet
var MyAllref Ref
//InitPacket 初始化
func InitPacket(content []byte) bool {

	err := xml.Unmarshal(content, &MyPacketGet)
	if err != nil {
		log.Fatal(err)
		return false
	}
	err = xml.Unmarshal(content, &MyAllref)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}
