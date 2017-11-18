//对json文件的操作

package gjson

import (
	"encoding/json"
)

//Js json数据结构体
type Js struct {
	mdata interface{}
}

//NewJSONString 通过string创建json结构
func NewJSONString(data string) (*Js, error) {
	j := new(Js)
	var f interface{}
	err := json.Unmarshal([]byte(data), &f)
	if err != nil {
		return j, err
	}
	j.mdata = f
	return j, err
}

//NewJSONByte 通过[]byte创建json结构
func NewJSONByte(data []byte) (*Js, error) {
	j := new(Js)
	var f interface{}
	err := json.Unmarshal(data, &f)
	if err != nil {
		return j, err
	}
	j.mdata = f
	return j, err
}

//NewGet 通过json结构创建json结构
func NewGet(j *Js, key string) *Js {
	value := new(Js)
	value.mdata = j.get(key)
	return value
}

//NewGetindex 通过json结构创建json结构
func NewGetindex(j *Js, i int) *Js {
	value := new(Js)
	value.mdata = j.getlist(i)
	return value
}

//Getnum 获取[]数组数量
func (j *Js) Getnum() int {
	if m, ok := (j.mdata).([]interface{}); ok {
		return len(m)
	}
	return 0
}

//IsValid 判断json结构是否有数据
func (j *Js) IsValid() bool {
	if nil == j.mdata {
		return false
	}
	return true
}

//Getuint16 通过key获取uint16
func (j *Js) Getuint16(key string) uint16 {
	data := j.get(key)
	if data != nil {
		if m, ok := data.(float64); ok {
			return uint16(m)
		}
	}
	return 0
}

//Getstring 通过key获取string
func (j *Js) Getstring(key string) string {
	data := j.get(key)
	if data != nil {
		if m, ok := data.(string); ok {
			return m
		}
	}
	return ""
}

//Getint32 通过key获取int32
func (j *Js) Getint32(key string) int32 {
	data := j.get(key)
	if data != nil {
		if m, ok := data.(float64); ok {
			return int32(m)
		}
	}
	return 0
}

//Getint 通过key获取int
func (j *Js) Getint(key string) int {
	data := j.get(key)
	if data != nil {
		if m, ok := data.(float64); ok {
			return int(m)
		}
	}
	return 0
}

//Getbool 通过key获取bool
func (j *Js) Getbool(key string) bool {
	data := j.get(key)
	if data != nil {
		if m, ok := data.(bool); ok {
			return m
		}
	}
	return false
}

func (j *Js) get(key string) interface{} {
	if m, ok := (j.mdata).(map[string]interface{}); ok {
		if data, oki := m[key]; oki {
			return data
		}
	}
	return nil
}
func (j *Js) getlist(i int) interface{} {

	if m, ok := (j.mdata).([]interface{}); ok {
		if i > 0 && i < len(m) {
			return m[i]
		}
	}
	return nil
}
