package business

import "sync"

// NowMonitorRse 当前监控信息
var NowMonitorRse sync.Map

func NowMonitorSet(now *NowMonitor) {
	NowMonitorRse.Store(now.HostId, now)
}

func NowMonitorGet(hostId string) *NowMonitor {
	if data, is := NowMonitorRse.Load(hostId); is {
		return data.(*NowMonitor)
	}
	return &NowMonitor{}
}

type NowMonitor struct {
	HostId     string
	Code       int
	Ms         int64
	ContrastMs int64
	PingMs     int64
}
