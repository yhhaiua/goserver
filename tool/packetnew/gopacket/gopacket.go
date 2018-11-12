package gopacket

import (
	"bytes"
	"fmt"
	"github.com/yhhaiua/goserver/tool/packetnew/structure"
	"os"
	"strconv"
	"strings"
)

const dir = "../../packets/"
//Create 创建go对应的数据包
func Create(path string)  {
	dirRef = path+dir +"objects"
	dirClient = path+dir +"client"
	dirServer = path+dir +"server"
	dirCode = path + dir
	os.MkdirAll(dirClient, os.ModePerm)
	os.MkdirAll(dirRef, os.ModePerm)
	os.MkdirAll(dirServer, os.ModePerm)

	conversionRef()
	conversionPacket()
	conversionCode()
}

var dirRef string
var dirClient string
var dirServer string
var dirCode string

var myCodeName []string
var myCodeID []int32


func conversionRef()  {

	for _, temref := range structure.MyAllref.MyRef {
		filename := strings.ToLower(temref.Name)
		filename = dirRef+"/" + filename + ".go"
		file, err := os.Create(filename)
		if err == nil {
			conversionOneRef(file, &temref)
			file.Close()
		}
	}
}
func conversionOneRef(file *os.File, data *structure.OneRef)  {

	var buf bytes.Buffer
	fmt.Fprint(&buf, "package objects\n\n")
	fmt.Fprintf(&buf, "//%s %s\n", data.Name, data.Des)
	fmt.Fprintf(&buf, "type %s struct {\n", data.Name)
	for _, tempfiled := range data.Field {
		conversionField(&tempfiled, &buf,false)
	}
	fmt.Fprint(&buf, "}\n\n")
	file.Write(buf.Bytes())
}

func  conversionField(field *structure.StField, tempbuf *bytes.Buffer,boobject bool)  {
	var fieldtype string
	switch field.Type {
	case "ref":
		if boobject{
			fieldtype = "objects."+field.RefType
		}else{
			fieldtype = field.RefType
		}

	case "refArray":
		if boobject{
			fieldtype = "[]objects."+field.RefType
		}else{
			fieldtype = "[]"+field.RefType
		}
	case "int":
		fieldtype = "int32"
	case "int[]":
		fieldtype = "[]int32"
	case "long":
		fieldtype = "int64"
	case "long[]":
		fieldtype = "[]int64"
	case "String":
		fieldtype = "string"
	case "String[]":
		fieldtype = "[]string"
	case "byte":
		fieldtype = "int8"
	case "byte[]":
		fieldtype = "[]int8"
	case "short":
		fieldtype = "int16"
	case "short[]":
		fieldtype = "[]int16"
	default:
		fieldtype = field.Type
	}

	fmt.Fprintf(tempbuf, "	%s %s //%s\n", capitalize(field.Name), fieldtype, field.Des)
}

func conversionPacket()  {

	for _, temppacket := range structure.MyPacketGet.Packet {

		filename := strings.ToLower(temppacket.Name)
		if temppacket.Type == "0" || temppacket.Type == "3"{
			filename = dirClient +"/" + filename + ".go"
			file, err := os.Create(filename)
			if err == nil {
				packetclient(file, &temppacket)
				cmd, _ := strconv.Atoi(temppacket.Id)
				myCodeName = append(myCodeName, temppacket.Name)
				myCodeID = append(myCodeID, int32(cmd))
				file.Close()
			}
		}else{
			filename = dirServer +"/" + filename + ".go"
			file, err := os.Create(filename)
			if err == nil {
				packetserver(file, &temppacket)
				cmd, _ := strconv.Atoi(temppacket.Id)
				myCodeName = append(myCodeName, temppacket.Name)
				myCodeID = append(myCodeID, int32(cmd))
				file.Close()
			}
		}

	}
}
func packetserver(file *os.File, data *structure.OnePacketGet)  {
	var buf bytes.Buffer
	fmt.Fprint(&buf, "package server\n\n")

	if(data.Field != nil){
		for _, tempfiled := range data.Field {
			if tempfiled.RefType != ""{
				fmt.Fprint(&buf, "import \"github.com/yhhaiua/goserver/packets/objects\"\n\n")
				break
			}
		}
	}
	fmt.Fprintf(&buf, "//%s %s\n", data.Name, data.Des)
	fmt.Fprintf(&buf, "type %s struct {\n	Code int32\n", data.Name)

	if(data.Field != nil){
		for _, tempfiled := range data.Field {
			conversionField(&tempfiled, &buf,true)
		}
	}

	fmt.Fprint(&buf, "}\n\n")
	fmt.Fprintf(&buf, "//Init %s初始化\n", data.Name)
	fmt.Fprintf(&buf, "func (pcmd *%s) Init() {\n   pcmd.Code = %s\n}", data.Name, data.Id)
	file.Write(buf.Bytes())
}
func packetclient(file *os.File, data *structure.OnePacketGet)  {

	var buf bytes.Buffer
	fmt.Fprint(&buf, "package client\n\n")

	if(data.Field != nil){
		for _, tempfiled := range data.Field {
			if tempfiled.RefType != ""{
				fmt.Fprint(&buf, "import \"github.com/yhhaiua/goserver/packets/objects\"\n\n")
				break
			}
		}
	}
	fmt.Fprintf(&buf, "//%s %s\n", data.Name, data.Des)
	fmt.Fprintf(&buf, "type %s struct {\n	Code int32\n", data.Name)

	if(data.Field != nil){
		for _, tempfiled := range data.Field {
			conversionField(&tempfiled, &buf,true)
		}
	}
	fmt.Fprint(&buf, "}\n\n")
	fmt.Fprintf(&buf, "//Init %s初始化\n", data.Name)
	fmt.Fprintf(&buf, "func (pcmd *%s) Init() {\n   pcmd.Code = %s\n}", data.Name, data.Id)
	file.Write(buf.Bytes())
}
// capitalize 字符首字母大写
func capitalize(str string) string {
	var upperStr string
	vv := []rune(str)   // 后文有介绍
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 {  // 后文有介绍
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				fmt.Println("Not begins with lowercase letter,")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}
func conversionCode() {
	filename := dirCode + "code.go"
	file, err := os.Create(filename)
	if err == nil {
		defer file.Close()
		var buf bytes.Buffer
		fmt.Fprint(&buf, "package packets\n\n")
		fmt.Fprint(&buf, "//包的key\nconst (\n")

		for i, temName := range myCodeName {
			fmt.Fprintf(&buf, "	%s = %d\n", temName, myCodeID[i])
		}
		fmt.Fprint(&buf, ")")
		file.Write(buf.Bytes())
	}
}
func clientPackets()  {
	//filename := dirCode + "clientpackets.go"
	//file, err := os.Create(filename)
	//if err == nil {
	//	defer file.Close()
	//	var buf bytes.Buffer
	//	fmt.Fprint(&buf, "package packets\n\n")
	//	fmt.Fprint(&buf, "import \"github.com/yhhaiua/goserver/packets/client\"\n\n")
	//	fmt.Fprint(&buf, "import \"github.com/yhhaiua/goserver/packets/server\"\n\n")
	//}
}