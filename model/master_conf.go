package model

// 平台设置 默认值

var (
	ContrastUriDefault  = "www.baidu.com"
	ContrastTimeDefault = 1000
	PingDefault         = "101.226.4.6"
	LogSaveDayDefault   = 7
)

type MasterConf struct {
	// 对照组
	ContrastUri  string
	ContrastTime int64
	// 用于检查当前网络
	Ping string
	// 日志保留天数
	LogSaveDay int
}

func GetMasterConf() *MasterConf {
	masterConf := &MasterConf{}
	err := DB.Get(MasterConfTable, MasterConfKey, masterConf)
	if err != nil {
		return &MasterConf{
			ContrastUri:  ContrastUriDefault,
			ContrastTime: int64(ContrastTimeDefault),
			Ping:         PingDefault,
			LogSaveDay:   LogSaveDayDefault,
		}
	}
	return masterConf
}
