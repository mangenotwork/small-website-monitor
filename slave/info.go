package slave

import (
	"fmt"
	"os"
	"runtime"
	"small-website-monitor/global"
	"small-website-monitor/model"
)

// 可扩展slave

// SlaveInfo 检测器信息
type SlaveInfo struct {
	Version         string // 版本
	HttpSubassembly string // 请求组件
	IP              string // IP
	Address         string // IP所在地址
	OSVersion       string // 系统版本
	//
}

func GetSlaveInfo() *SlaveInfo {
	myIP := model.GetMyIP()
	return &SlaveInfo{
		Version:         global.Version,
		HttpSubassembly: global.HttpSubassembly,
		IP:              myIP.IP,
		Address:         myIP.Address,
		OSVersion:       GetOSInfo(),
	}
}

func GetOSInfo() string {
	hostName, _ := os.Hostname()
	return fmt.Sprintf("主机名称:%s;系统类型:%s;系统架构:%s(%d核);",
		hostName, runtime.GOOS, runtime.GOARCH, runtime.GOMAXPROCS(0))
}
