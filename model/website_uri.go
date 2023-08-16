package model

import (
	"github.com/mangenotwork/common/log"
	gt "github.com/mangenotwork/gathertool"
)

// WebSiteUri 站点的Uri存储
type WebSiteUri struct {
	HostID  string
	AllUri  []string
	ExtLink []string // 外链
	BadLink []string
}

func NewWebSiteUri(hostId string) *WebSiteUri {
	return &WebSiteUri{
		HostID:  hostId,
		AllUri:  make([]string, 0),
		ExtLink: make([]string, 0),
		BadLink: make([]string, 0),
	}
}

func (m *WebSiteUri) Add() error {
	return DB.Set(WebSiteURITable, m.HostID, m)
}

func (m *WebSiteUri) Get() (*WebSiteUri, error) {
	err := DB.Get(WebSiteURITable, m.HostID, m)
	return m, err
}

func (m *WebSiteUri) Collect(rootURI string, depth int) {
	gt.ApplicationTerminalOut = false
	hostScan := gt.NewHostScanUrl(rootURI, depth)
	m.AllUri, _ = hostScan.Run()
	//log.Info("m.AllUri = ", m.AllUri)
	extLinks := gt.NewHostScanExtLinks(rootURI)
	m.ExtLink, _ = extLinks.Run()
	//log.Info("m.ExtLink = ", m.ExtLink)
	// 死链接监测
	var badNum = 0
	badLinks := gt.NewHostScanBadLink(rootURI, depth)
	m.BadLink, badNum = badLinks.Run()
	//log.Info("m.BadLink = ", m.BadLink)
	if badNum > 0 {
		// TODO 发送邮件通知存在死链接
		log.Error("存在死链接")
	}
	err := m.Add()
	if err != nil {
		log.Error("保存数据失败, ", err.Error())
	}
	return
}
