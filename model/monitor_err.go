package model

// MonitorErrInfo 监测器错误信息
type MonitorErrInfo struct {
	List []string
}

func NewMonitorErrInfo() *MonitorErrInfo {
	return &MonitorErrInfo{
		List: make([]string, 0),
	}
}
