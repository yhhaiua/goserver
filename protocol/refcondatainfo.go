package protocol

//RefConDataInfo 连接数据结构
type RefConDataInfo struct {
	Svrid   int32  //服务器id
	Svrtype int32  //服务器类型
	Sip     string //服务器ip
	Sport   string //服务器端口
}
