package comsvrsrc

//服务器类型
const (
	SERVERTYPELOGIN  = 1000 // 登录服务器
	SERVERTYPEDUTY   = 1100 // 登录职守
	SERVERTYPEGATE   = 2000 // 网关服务器
	SERVERTYPEGAME   = 3000 // 游戏服务器
	SERVERTYPEMANAGE = 4000 // 管理服务器
)

const (
	//CHECKDATACODE 服务器间数据检测
	CHECKDATACODE = 0x55884433
)
