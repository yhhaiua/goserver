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
)

type StField struct {
	Define    string `xml:"define,attr"`
	Field    string `xml:"field,attr"`
	Enum string `xml:"enum,attr"`
}
type Configs struct {
	Fields []StField `xml:"field"`
}
type AllConfig struct {
	MyConfig Configs `xml:"config"`
}

type  UnitOne struct {
	Id string  `xml:"id,attr"`
}
type UnitConfig struct {
	Unit []UnitOne `xml:"unit"`
}
var Allid map[int]int
var Reward map[int]int
var MyDir string

var AllData []string
func main() {
	args := os.Args
	if args == nil || len(args) < 2 {
		return
	}
	MyDir = args[1]

	Allid = make(map[int]int)
	Reward = make(map[int]int)
	readFile(MyDir+"card.xml")
	readFile(MyDir+"item.xml")
	readFile(MyDir+"expert.xml")
	readFile(MyDir+"cardchip.xml")
	readFile2(MyDir+"reward.xml")
	log.Printf("All num:%d",len(Allid))
	log.Printf("reward num:%d",len(Reward))

	info := dirents(MyDir)
	if info != nil {
		for _,oneinfo:= range info {
			if !oneinfo.IsDir() && oneinfo.Name() != "version.txt"{
				dir := fmt.Sprintf("%s%s",
					MyDir,
					oneinfo.Name())
				xmlfile(dir,oneinfo.Name())
			}

		}
	}
	if AllData == nil{
		log.Printf("SUCCESS")
		data:= fmt.Sprintf("SUCCESS")
		AllData = append(AllData,data)
	}
	if AllData != nil{
		createdir()
	}
}
func createdir() {
	filename := "./check.log"
	file, err := os.Create(filename)
	if err == nil {
		defer file.Close()
		var buf bytes.Buffer

		for _, temName := range AllData {
			fmt.Fprintf(&buf, "	%s\n", temName)
		}
		file.Write(buf.Bytes())
	}
}
func readFile(dir string)  {

	content, err := ioutil.ReadFile(dir)
	if err != nil {
		log.Fatal(err)
		return
	}
	readitemid(string(content))

}
func readitemid(content string)  {
	info := strings.Split(content,"\n")
	for _,onefo := range info{
		if strings.Contains(onefo,"<unit id="){
			bonus := strings.Split(onefo," ")
			for _,onebonus := range bonus {
				if strings.Contains(onebonus,"id="){

					nexstr := strings.Split(onebonus,"\"")
					if nexstr != nil && len(nexstr) >= 2{
						if nexstr[0] == "id="{
							id,_:=strconv.Atoi(nexstr[1])
							Allid[id] = id
						}
					}
				}
			}
		}
	}
}
func readFile2(dir string)  {

	content, err := ioutil.ReadFile(dir)
	if err != nil {
		log.Fatal(err)
		return
	}
	readitemid2(string(content))

}
func readitemid2(content string)  {
	info := strings.Split(content,"\n")
	for _,onefo := range info{
		if strings.Contains(onefo,"<unit id="){
			bonus := strings.Split(onefo," ")
			for _,onebonus := range bonus {
				if strings.Contains(onebonus,"id="){

					nexstr := strings.Split(onebonus,"\"")
					if nexstr != nil && len(nexstr) >= 2{
						if nexstr[0] == "id="{
							id,_:=strconv.Atoi(nexstr[1])
							Reward[id] = id
						}
					}
				}
			}
		}
	}
	for _,onefo := range info{
		if strings.Contains(onefo,"reward="){
			bonus := strings.Split(onefo," ")
			for _,onebonus := range bonus {
				if strings.Contains(onebonus,"reward="){
					nexstr := strings.Split(onebonus,"\"")
					if nexstr != nil && len(nexstr) >= 2{
						if nexstr[0] == "reward="{

							onestr := strings.Split(nexstr[1],",")
							for _,one := range onestr{
								ids := strings.Split(one,":")
								//id,_:=strconv.Atoi(ids[0])
								//log.Println(id)
								if ids !=nil{
									if ids[0] != "REWARD"{
										id,_:=strconv.Atoi(ids[0])
										_,ok := Allid[id]
										if !ok{
											log.Printf("reward.xml no item:%d",id)

											data:= fmt.Sprintf("reward.xml中有道具:%d不存在",id)
											AllData = append(AllData,data)
										}
									}else{
										id,_:=strconv.Atoi(ids[1])
										_,ok := Reward[id]
										if !ok{
											log.Printf("reward.xml no reward.xml id:%d",id)
											data:= fmt.Sprintf("表:reward.xml 中有奖励reward.xml:%d不存在",id)
											AllData = append(AllData,data)
										}
									}

								}
							}
						}
					}
				}
			}
		}
	}
}
func xmlfile(dir string,name string)  {
	//log.Println(dir)
	content, err := ioutil.ReadFile(dir)
	if err != nil {
		log.Fatal(err)
		return
	}

	var temp AllConfig
	err = xml.Unmarshal(content, &temp)
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

	if tempunit.Unit != nil{
		tempmap := make(map[string]bool)
		for _,info := range tempunit.Unit{
			_,ok :=tempmap[info.Id]
			if !ok {
				tempmap[info.Id]=true
			}else{
				log.Printf("%s have id the same:%s",name,info.Id)
				data:= fmt.Sprintf("表:%s中有id相同:%s",name,info.Id)
				AllData = append(AllData,data)
			}
		}
	}
	define := findbonus(temp.MyConfig)
	if define != nil{
		readxml(string(content),name,define)
	}
}
func contains(str string,define []string) bool  {
	for _,info := range define{
		if strings.Contains(str,info){
			return true
		}
	}
	return false
}
func check(str string,define []string) bool {
	for _,info := range define{
		temp:= info +"="
		if temp == str{
			return true
		}
	}
	return  false
}
func readxml(content ,dir string,define []string)  {
	info := strings.Split(content,"\n")
	for _,onefo := range info{
		if contains(onefo,define) && !strings.Contains(onefo,"develop.game.common.convertor.BonusConvertor") && !strings.Contains(onefo,"<!--") {
			bonus := strings.Split(onefo," ")
			for _,onebonus := range bonus {
				if contains(onebonus,define){
					nexstr := strings.Split(onebonus,"\"")
					if nexstr != nil && len(nexstr) >= 2{
						if check(nexstr[0],define){
							onestr := strings.Split(nexstr[1],",")
							for _,one := range onestr{
								ids := strings.Split(one,":")
								//id,_:=strconv.Atoi(ids[0])
								//log.Println(id)
								if ids !=nil{
									if ids[0] != "REWARD"{
										id,_:=strconv.Atoi(ids[0])
										_,ok := Allid[id]
										if !ok{
											log.Printf("%s no item :%d",dir,id)
											data:= fmt.Sprintf("表:%s中有道具:%d不存在",dir,id)
											AllData = append(AllData,data)
										}
									}else{
										id,_:=strconv.Atoi(ids[1])
										_,ok := Reward[id]
										if !ok{
											log.Printf("%s no reward.xml id:%d",dir,id)
											data:= fmt.Sprintf("表:%s中有奖励reward.xml:%d不存在",dir,id)
											AllData = append(AllData,data)
										}
									}

								}
							}
						}else{
							//log.Println(nexstr)
						}

					}else{
						log.Println(nexstr)
						data:= fmt.Sprintf("表:%s出现错误:%s",dir,nexstr)
						AllData = append(AllData,data)
					}

				}
			}
		}
	}
}
func findbonus(temp Configs) []string  {
	var tempstr []string
	for _,info:= range temp.Fields{
		if info.Enum == "develop.game.common.convertor.BonusConvertor"{
			tempstr = append(tempstr,info.Define)
		}
	}
	return tempstr
}
func dirents(dir string) []os.FileInfo {

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return entries
}