package main

import (
	"fmt"
	"github.com/yhhaiua/goserver/tool/packetnew/javapacket"
	"github.com/yhhaiua/goserver/tool/packetnew/structure"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)


func getCurrentPath() string {
	s, err := exec.LookPath(os.Args[0])
	checkErr(err)
	i := strings.LastIndex(s, "\\")
	path := string(s[0 : i+1])
	return path
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {

	path := getCurrentPath()
	fmt.Println(path)

	content, err := ioutil.ReadFile(path+"packet.xml")
	if err != nil {
		log.Fatal(err)
		return
	}
	//初始化
	if structure.InitPacket(content){

		//gopacket.Create(path)
		javapacket.Create(path)
	}

}

