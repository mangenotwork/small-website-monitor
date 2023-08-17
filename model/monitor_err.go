package model

import (
	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"
)

// MonitorErrInfo 监测器错误信息
type MonitorErrInfo struct {
	List []string
}

func NewMonitorErrInfo() *MonitorErrInfo {
	return &MonitorErrInfo{
		List: make([]string, 0),
	}
}

func (m *MonitorErrInfo) Get() error {
	return DB.Get(MonitorErrInfoTable, MonitorErrInfoKey, m)
}

func (m *MonitorErrInfo) Add(errStr string) error {
	err := m.Get()
	if err != nil && err != ISNULL {
		return err
	}
	m.List = append(m.List, errStr)
	return DB.Set(MonitorErrInfoTable, MonitorErrInfoKey, m)
}

func (m *MonitorErrInfo) Clear() error {
	m.List = m.List[:0]
	return DB.Set(MonitorErrInfoTable, MonitorErrInfoKey, m)
}

func SetMonitorErrInfo(errStr string) {
	errInfo := NewMonitorErrInfo()
	addErr := errInfo.Add(utils.NowDate() + errStr)
	if addErr != nil {
		log.Error(addErr)
	}
}
