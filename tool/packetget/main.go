package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/yhhaiua/goserver/common/gpacket"
)

const dir = "../../protocol/"

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
var myCodeMap map[string]uint16
var myref Ref

func main() {

	myCodeMap = make(map[string]uint16)

	content, err := ioutil.ReadFile("packet.xml")
	if err != nil {
		log.Fatal(err)
	}

	err = xml.Unmarshal(content, &myPacketGet)
	if err != nil {
		log.Fatal(err)
		return
	}
	//log.Println(myPacketGet)

	err = xml.Unmarshal(content, &myref)
	if err != nil {
		log.Fatal(err)
		return
	}
	//log.Println(myref)

	conversionGo()
}

func conversionGo() {

	for _, temref := range myref.MyRef {
		filename := strings.ToLower(temref.Name)
		filename = dir + filename + ".go"
		file, err := os.Create(filename)
		if err == nil {
			oneconversionref(file, &temref)
			file.Close()
		}
	}
	for _, temppacket := range myPacketGet.Packet {

		filename := strings.ToLower(temppacket.Name)
		filename = dir + filename + ".go"
		file, err := os.Create(filename)
		if err == nil {
			oneconversion(file, &temppacket)
			icmd, _ := strconv.Atoi(temppacket.Cmd)
			isupcmd, _ := strconv.Atoi(temppacket.Supcmd)
			myCodeMap[temppacket.Name] = gpacket.GetValue(uint8(icmd), uint8(isupcmd))
			file.Close()
		}
	}
	codeconversion()
}

func oneconversionref(file *os.File, data *OneRef) {
	var buf bytes.Buffer
	fmt.Fprint(&buf, "package protocol\n\n")
	fmt.Fprintf(&buf, "//%s %s\n", data.Name, data.Des)
	fmt.Fprintf(&buf, "type %s struct {\n", data.Name)
	for _, tempfiled := range data.Field {
		filedconversion(&tempfiled, &buf)
	}
	fmt.Fprint(&buf, "}\n\n")
	file.Write(buf.Bytes())
}
func oneconversion(file *os.File, data *OnePacketGet) {

	var buf bytes.Buffer
	fmt.Fprint(&buf, "package protocol\n\nimport \"github.com/yhhaiua/goserver/common/gpacket\"\n\n")
	fmt.Fprintf(&buf, "//%s %s\n", data.Name, data.Des)
	fmt.Fprintf(&buf, "type %s struct {\n	gpacket.BaseCmd\n", data.Name)
	for _, tempfiled := range data.Field {
		filedconversion(&tempfiled, &buf)
	}
	fmt.Fprint(&buf, "}\n\n")
	fmt.Fprintf(&buf, "//Init %s初始化\n", data.Name)
	fmt.Fprintf(&buf, "func (pcmd *%s) Init() {\n   pcmd.Cmd = %s\n	  pcmd.SupCmd = %s\n}", data.Name, data.Cmd, data.Supcmd)
	file.Write(buf.Bytes())
}

func filedconversion(field *StField, tempbuf *bytes.Buffer) {
	switch field.Type {
	case "ref":
		fmt.Fprintf(tempbuf, "	%s %s\n", field.Name, field.RefType)
	case "refArray":
		fmt.Fprintf(tempbuf, "	%s []%s\n", field.Name, field.RefType)
	default:
		fmt.Fprintf(tempbuf, "	%s %s\n", field.Name, field.Type)
	}
}

func codeconversion() {
	filename := dir + "code.go"
	file, err := os.Create(filename)
	if err == nil {
		defer file.Close()
		var buf bytes.Buffer
		fmt.Fprint(&buf, "package protocol\n\n")
		fmt.Fprint(&buf, "//包的key\nconst (\n")
		for key, value := range myCodeMap {
			fmt.Fprintf(&buf, "	%sCode = %d\n", key, value)
		}
		fmt.Fprint(&buf, ")")
		file.Write(buf.Bytes())
	}
}
