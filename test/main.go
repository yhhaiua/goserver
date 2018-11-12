package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"math"
	"net/url"
	"strings"
)

const (
	width, height = 600, 320            // canvas size in pixels
	cells         = 100                 // number of grid cells
	xyrange       = 30.0                // axis ranges (-xyrange..+xyrange)
	xyscale       = width / 2 / xyrange // pixels per x or y unit
	zscale        = height * 0.4        // pixels per z unit
	angle         = math.Pi / 6         // angle of x, y axes (=30°)
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)

func main() {

	//https://api.urlshare.cn/v3/user/get_info?openid=12345&openkey=12345&pf=qzone&appid=1105583577&format=json&userip=10.0.0.1&test=%2A&sig=upLbDHoONj4UEknDSvcun1yEqnk%3D
	//https://api.urlshare.cn/v3/user/get_info?pf=qzone&userip=10.0.0.1&test=%2A&sig=upLbDHoONj4UEknDSvcun1yEqnk%3D

	//sigCreate("/v3/user/get_info","appid=1105583577&format=json&openid=12345&openkey=12345&pf=qzone&test=%2A&userip=10.0.0.1","228bf094169a40a3bd188ba37ebe8723")
	sigCreate("/v3/user/get_info","appid=123456&format=json&openid=11111111111111111&openkey=2222222222222222&pf=qzone&userip=112.90.139.30","228bf094169a40a3bd188ba37ebe8723")
}

func sigCreate(urlinit,urlstr,openkey string)  {
	str0 :="GET"
	str1:= encodeURIComponent(urlinit)
	fmt.Println(str1)
	str2 := encodeURIComponent(urlstr)
	fmt.Println(str2)
	str3 := str0+"&"+str1+"&"+str2

	key :=  []byte(openkey+"&")
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(str3))
	uEnc := base64.URLEncoding.EncodeToString(mac.Sum(nil))
	fmt.Println(uEnc)
	fmt.Println(encodeURIComponent(uEnc))
}
func encodeURIComponent(str string) string {
	r := url.QueryEscape(str)
	r = strings.Replace(r, "+", "%20", -1)
	return r
}
func sign(urlStr string)  {
	l, err := url.ParseQuery(urlStr)
	fmt.Println(l, err)
	l2, err2 := url.ParseRequestURI(urlStr)
	fmt.Println(l2, err2)

	l3, err3 := url.Parse(urlStr)
	fmt.Println("1")
	fmt.Println(l3, err3)
	fmt.Println("2")
	fmt.Println(l3.Query())
	fmt.Println("3")
	fmt.Println(l3.Query().Encode())
}
func corner(i, j int) (float64, float64) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	z := f(x, y)

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	return math.Sin(r) / r
}
