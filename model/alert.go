package model

// WebSiteAlert 监控报警信息
type WebSiteAlert struct {
	HostID string
	List   []*AlertData
}

type AlertData struct {
	Date          string // 监测的时间
	Uri           string // 出现问题的URI
	UriCode       int    // URI响应码
	UriMs         int64  // URI响应时间
	ContrastUriMs int64  // 对照组URI响应时间
	PingMs        int64  // ping 当前网络情况
	Msg           string // 报警信息
}

func NewWebSiteAlert(hostId string) *WebSiteAlert {
	return &WebSiteAlert{
		HostID: hostId,
		List:   make([]*AlertData, 0),
	}
}

func (m *WebSiteAlert) Get() error {
	return DB.Get(WebSiteAlertTable, m.HostID, m)
}

func (m *WebSiteAlert) Add(alert *AlertData) error {
	err := m.Get()
	if err != nil && err != ISNULL {
		return err
	}
	m.List = append(m.List, alert)
	return DB.Set(WebSiteAlertTable, m.HostID, m)
}

func (m *WebSiteAlert) Clear() error {
	m.List = m.List[:0]
	return DB.Set(WebSiteAlertTable, m.HostID, m)
}
