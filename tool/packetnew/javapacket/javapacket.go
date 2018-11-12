package javapacket

import (
	"bytes"
	"fmt"
	"github.com/yhhaiua/goserver/tool/packetnew/structure"
	"os"
	"strconv"
)

const dir = "../../packets/"
//Create 创建java对应的包
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
	//conversionCode()
}


var dirRef string
var dirClient string
var dirServer string
var dirCode string

var myCodeName []string
var myCodeID []int32


func conversionRef()  {

	for _, temref := range structure.MyAllref.MyRef {
		filename := temref.Name
		filename = dirRef+"/" + filename + ".java"
		file, err := os.Create(filename)
		if err == nil {
			conversionOneRef(file, &temref)
			file.Close()
		}
	}
}

func conversionOneRef(file *os.File, data *structure.OneRef)  {

	var buf bytes.Buffer
	fmt.Fprint(&buf, "package develop.packets.objects;\nimport engine.net.*;\n")
	fmt.Fprintf(&buf, "//%s\n", data.Des)
	fmt.Fprintf(&buf, "public class %s extends CValue\n{\n", data.Name)
	for _, tempfiled := range data.Field {
		conversionField(&tempfiled, &buf)
	}
	fmt.Fprint(&buf, "	public void read(engine.net.NativeBuffer buf)\n	{\n")
	fmt.Fprint(&buf, "	}\n")
	fmt.Fprint(&buf, "}\n")
	file.Write(buf.Bytes())
}

func  conversionField(field *structure.StField, tempbuf *bytes.Buffer)  {
	var fieldtype string
	switch field.Type {
	case "ref":
		fieldtype = field.RefType

	case "refArray":
		fieldtype = field.RefType+"[]"
	case "string":
		fieldtype = "String"
	case "string[]":
		fieldtype = "String[]"
	default:
		fieldtype = field.Type
	}

	fmt.Fprintf(tempbuf, "	public %s %s //%s\n",fieldtype,field.Name, field.Des)
}

func  conversionFieldValue(field *structure.StField, tempbuf *bytes.Buffer)  {
	var value string
	switch field.Type {
	case "int":
		value = field.Name+"=buf.readInt();"
	case "short":
		value = field.Name+"=buf.readShort();"
	case "byte":
		value = field.Name+"=buf.readByte();"
	case "long":
		value = field.Name+"=buf.readLong();"
	case "float":
		value = field.Name+"=buf.readFloat();"
	case "boolean":
		value = field.Name+"=buf.readBoolean();"
	case "String","string":
		value = field.Name+"=buf.readUTF();"
	case "int[]":
		value = field.Name+"=buf.readIntArray();"
	case "short[]":
		value = field.Name+"=buf.readShortArray();"
	case "byte[]":
		value = field.Name+"=buf.readByteArray();"
	case "long[]":
		value = field.Name+"=buf.readLongArray();"
	case "float[]":
		value = field.Name+"=buf.readFloatArray();"
	case "String[]","string[]":
		value = field.Name+"=buf.readUTFArray();"
	case "int[][]":
		value = field.Name+"=buf.readIntTwoArray();"
	default:

	}
	fmt.Fprintf(tempbuf, "		public %s\n",value)
}
func conversionPacket()  {

	for _, temppacket := range structure.MyPacketGet.Packet {

		filename := temppacket.Name
		if temppacket.Type == "0" || temppacket.Type == "3"{
			filename = dirClient +"/" + filename + ".java"
			file, err := os.Create(filename)
			if err == nil {
				packetclient(file, &temppacket)
				cmd, _ := strconv.Atoi(temppacket.Id)
				myCodeName = append(myCodeName, temppacket.Name)
				myCodeID = append(myCodeID, int32(cmd))
				file.Close()
			}
		}else{
			filename = dirServer +"/" + filename + ".java"
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
	fmt.Fprint(&buf, "package develop.packets.server;\nimport engine.net.*;\n")

	if(data.Field != nil){
		for _, tempfiled := range data.Field {
			if tempfiled.RefType != ""{
				fmt.Fprint(&buf, "import develop.packets.objects.*;\n")
				break
			}
		}
	}
	fmt.Fprintf(&buf, "//%s\n", data.Des)
	fmt.Fprintf(&buf, "public class %s extends CPacket\n{\n", data.Name)

	if(data.Field != nil){
		for _, tempfiled := range data.Field {
			conversionField(&tempfiled, &buf)
		}
	}
	fmt.Fprint(&buf, "}\n")
	file.Write(buf.Bytes())
}
func packetclient(file *os.File, data *structure.OnePacketGet)  {

	var buf bytes.Buffer
	fmt.Fprint(&buf, "package develop.packets.client;\nimport engine.net.*;\n")

	if(data.Field != nil){
		for _, tempfiled := range data.Field {
			if tempfiled.RefType != ""{
				fmt.Fprint(&buf, "import develop.packets.objects.*;\n")
				break
			}
		}
	}
	fmt.Fprintf(&buf, "//%s\n", data.Des)
	fmt.Fprintf(&buf, "public class %s extends CPacket\n{\n", data.Name)

	if(data.Field != nil){
		for _, tempfiled := range data.Field {
			conversionField(&tempfiled, &buf)
		}
	}
	fmt.Fprint(&buf, "}\n")
	file.Write(buf.Bytes())
}