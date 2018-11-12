package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type UnitOne struct {
	Id   string `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

type JsonUnitOne struct {
	Id   string `json:"itemId"`
	Name string `json:"itemName"`
}
type UnitConfig struct {
	Unit []UnitOne `xml:"unit"`
}

var allconfig []JsonUnitOne

func main() {
	args := os.Args
	if args == nil || len(args) < 2 {
		return
	}
	MyDir := args[1]
	//MyDir := "D:/GoProject/src/github.com/yhhaiua/goserver/tool/itemtojson/"
	xmlfile(MyDir, "item.xml")
	xmlfile(MyDir, "card.xml")
	xmlfile(MyDir, "cardchip.xml")
	xmlfile(MyDir, "expert.xml")

	if allconfig != nil {

		filename := "./item.xml"
		file, err := os.Create(filename)
		if err == nil {
			defer file.Close()
			var buf bytes.Buffer
			fmt.Fprint(&buf, "[\n")
			vlen := len(allconfig)
			i := 0
			for _,temp := range allconfig {
				i++
				data, err := json.Marshal(temp)
				if err == nil {
					fmt.Fprint(&buf,string(data))
					if(i < vlen){
						fmt.Fprint(&buf,",")
					}
					fmt.Fprint(&buf,"\n")
				}else{
					log.Println("错误json解析")
				}
			}
			fmt.Fprint(&buf, "]\n")
			file.Write(buf.Bytes())
		}

		log.Println("success")
	}

}
func xmlfile(dir string, name string) {
	content, err := ioutil.ReadFile(dir + name)
	if err != nil {
		log.Fatal(err)
		return
	}
	var tempunit UnitConfig
	err = xml.Unmarshal(content, &tempunit)
	if err != nil {
		log.Fatal(err)
		return
	}
	if tempunit.Unit != nil {
		for _, temp := range tempunit.Unit {
			var jsondata JsonUnitOne
			jsondata.Id = temp.Id
			jsondata.Name = temp.Name
			allconfig = append(allconfig, jsondata)
		}
	}
}
