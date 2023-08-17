package model

import (
	"fmt"
	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"
)

// WebSitePoint 站点监测点
type WebSitePoint struct {
	HostID string
	Uri    []string
}

func NewWebSitePoint(hostId string) *WebSitePoint {
	return &WebSitePoint{
		HostID: hostId,
		Uri:    make([]string, 0),
	}
}

func (m *WebSitePoint) Add(uri string) error {
	err := m.Get()
	if err != nil && err != ISNULL {
		return err
	}
	log.Info("m.Uri = ", m.Uri)
	for _, v := range m.Uri {
		if v == uri {
			return fmt.Errorf("监测点存在!")
		}
	}
	m.Uri = append(m.Uri, uri)
	log.Info("m.Uri = ", m.Uri)
	return DB.Set(WebSitePointTable, m.HostID, m)
}

func (m *WebSitePoint) Get() error {
	return DB.Get(WebSitePointTable, m.HostID, m)
}

func (m *WebSitePoint) Del(uri string) error {
	err := m.Get()
	if err != nil && err != ISNULL {
		return err
	}
	for n, v := range m.Uri {
		if v == uri {
			log.Info("n = ", n)
			m.Uri = append(m.Uri[:n], m.Uri[n+1:]...)
			break
		}
	}
	return DB.Set(WebSitePointTable, m.HostID, m)
}

func (m *WebSitePoint) Random() string {
	err := m.Get()
	if err != nil {
		log.Error(err)
		return ""
	}
	return utils.RandomString(m.Uri)
}
