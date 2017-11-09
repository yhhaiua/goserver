package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const dir = "./"

//StField field
type StField struct {
	Name    string `xml:"name,attr"`
	Type    string `xml:"type,attr"`
	RefType string `xml:"refType,attr"`
}

//OnePacketGet 结构
type OnePacketGet struct {
	Cmd    string    `xml:"cmd,attr"`
	Supcmd string    `xml:"supcmd,attr"`
	Name   string    `xml:"name,attr"`
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

var myPacketGet PacketGet
var myRefMap map[string]OneRef

func main() {

	myRefMap = make(map[string]OneRef)

	content, err := ioutil.ReadFile("packet.xml")
	if err != nil {
		log.Fatal(err)
	}

	err = xml.Unmarshal(content, &myPacketGet)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println(myPacketGet)

	var myref Ref
	err = xml.Unmarshal(content, &myref)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println(myref)

	for _, temref := range myref.MyRef {
		myRefMap[temref.Name] = temref
	}
	conversionGo()
}

func conversionGo() {

	for _, temppacket := range myPacketGet.Packet {

		filename := strings.ToLower(temppacket.Name)
		filename = dir + filename + ".go"
		file, err := os.Create(filename)
		if err == nil {
			defer file.Close()
			oneconversion(file, &temppacket)
		}
	}
}

func oneconversion(file *os.File, data *OnePacketGet) {

	var buf bytes.Buffer
	fmt.Fprint(&buf, "package protocol\n\nimport \"github.com/yhhaiua/goserver/common/gpacket\"\n\n")
	fmt.Fprintf(&buf, "//%s %s\n", data.Name, data.Des)
	fmt.Fprintf(&buf, "type %s struct {\n	gpacket.BaseCmd\n", data.Name)
	fmt.Fprint(&buf, "}\n\n")
	file.Write(buf.Bytes())
}
