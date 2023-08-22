package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/common/conf"
	"github.com/mangenotwork/common/ginHelper"
	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"
	gt "github.com/mangenotwork/gathertool"
	"net/http"
	"small-website-monitor/business"
	"small-website-monitor/global"
	"small-website-monitor/model"
	"sort"
	"strings"
	"time"
)

func Login(c *gin.Context) {
	user := c.PostForm("user")
	password := c.PostForm("password")
	for _, v := range conf.Conf.Default.User {
		if user == v.Name && password == v.Password {
			j := utils.NewJWT(conf.Conf.Default.Jwt.Secret, conf.Conf.Default.Jwt.Expire)
			j.AddClaims("name", user)
			token, tokenErr := j.Token()
			if tokenErr != nil {
				log.Error("生产token错误， err = ", tokenErr)
			}
			c.SetCookie(global.UserToken, token, global.TokenExpires,
				"/", "", false, true)
			c.Redirect(http.StatusFound, "/home")
			return
		}
	}
	c.HTML(200, "err.html", gin.H{
		"Title": conf.Conf.Default.App.Name,
		"err":   "账号或密码错误",
	})
	return
}

func Out(c *gin.Context) {
	c.SetCookie("sign", "", 60*60*24*7, "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}

type WebsiteAddParam struct {
	Host           string `json:"host"`
	HealthUri      string `json:"healthUri"`
	Rate           int64  `json:"rate"`
	AlarmResTime   int64  `json:"alarmResTime"`
	UriDepth       int64  `json:"uriDepth"`
	UriUpdateRate  int64  `json:"uriUpdateRate"`
	AlarmStateCode string `json:"alarmStateCode"`
}

func WebsiteAdd(c *ginHelper.GinCtx) {
	param := &WebsiteAddParam{}
	err := c.GetPostArgs(param)
	if err != nil {
		c.APIOutPutError(err, "参数或参数类型错误")
		return
	}
	if len(param.Host) < 1 {
		c.APIOutPutError(nil, "参数错误: host不能为空")
		return
	}
	ctx, err := gt.Get(param.HealthUri)
	if err != nil || !business.HttpCodeIsHealth(ctx.StateCode) {
		c.APIOutPutError(nil, fmt.Sprintf("生命监测uri网络不可达,http state code :%d", ctx.StateCode))
		return
	}
	if param.Rate < 1 {
		param.Rate = 1
	}
	if param.UriDepth < 1 {
		param.UriDepth = 1
	}
	if param.UriUpdateRate < 1 {
		param.UriUpdateRate = 1
	}
	// 小于100ms的响应时间不应该被报警
	if param.AlarmResTime < 100 {
		param.AlarmResTime = 100
	}
	website := &model.WebSite{
		Host:          param.Host,
		HealthUri:     utils.URIStr(param.HealthUri),
		Rate:          param.Rate,
		UriDepth:      param.UriDepth,
		UriUpdateRate: param.UriUpdateRate,
		AlarmResTime:  param.AlarmResTime,
		HostIP:        ctx.Req.RemoteAddr,
		Created:       time.Now().Unix(),
	}
	websiteId, err := website.Add()
	if err != nil {
		c.APIOutPutError(err, "保存数据失败")
		return
	}
	go func() {
		// 更新站点对象
		business.Push()
		// 采集站点URI相关信息
		webSiteUri := model.NewWebSiteUri(websiteId)
		webSiteUri.Collect(website.HealthUri, int(website.UriDepth))
	}()
	c.APIOutPut("", "添加成功")
	return
}

func CaseT(c *ginHelper.GinCtx) {
	c.APIOutPut("test", "")
}

type WebsiteListOut struct {
	List     []*WebsiteListData
	Count    int
	PageList []*ginHelper.Page
}

type WebsiteListData struct {
	*model.WebSite
	*business.NowMonitor
	AlertCount int
}

func WebsiteList(c *ginHelper.GinCtx) {
	pg := c.GetQueryInt("pg")
	if pg < 1 {
		pg = 1
	}
	size := 10
	websiteListData := make([]*WebsiteListData, 0)
	websiteList, count, err := new(model.WebSite).List(pg, size)
	if err != nil {
		c.APIOutPutError(err, "获取失败")
		return
	}
	for _, v := range websiteList {
		nowRse := business.NowMonitorGet(v.ID)
		alert := model.NewWebSiteAlert(v.ID)
		alertCount, err := alert.Len()
		if err != nil && err != model.ISNULL {
			log.Error(err)
		}
		websiteListData = append(websiteListData, &WebsiteListData{
			v, nowRse, alertCount,
		})
	}
	c.APIOutPut(&WebsiteListOut{
		List:     websiteListData,
		Count:    count,
		PageList: c.PageList(pg, 5, count, size, ""),
	}, "")
	return
}

func MailInit(c *ginHelper.GinCtx) {
	data := model.IsMail()
	c.APIOutPut(data, "")
	return
}

type MailConfParam struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	From     string `json:"from"`
	AuthCode string `json:"authCode"`
	ToList   string `json:"toList"`
}

func mailSet(c *ginHelper.GinCtx) error {
	param := &MailConfParam{}
	err := c.GetPostArgs(param)
	if err != nil {
		return err
	}
	if len(param.From) < 1 {
		return fmt.Errorf("%s", "发件人不能为空!")
	}
	if len(param.Host) < 1 {
		return fmt.Errorf("%s", "邮件服务器不能为空!")
	}
	if len(param.AuthCode) < 1 {
		return fmt.Errorf("%s", "邮件服务授权码不能为空!")
	}
	if len(param.ToList) < 1 {
		return fmt.Errorf("%s", "通知收件人不能为空!")
	}
	param.ToList = utils.CleaningStr(param.ToList)
	toList := strings.Split(param.ToList, ";")
	m := &model.MailData{
		From:     param.From,
		AuthCode: param.AuthCode,
		Host:     param.Host,
		Port:     param.Port,
		ToList:   toList,
	}
	return model.SetMail(m)
}

func MailConf(c *ginHelper.GinCtx) {
	err := mailSet(c)
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	c.APIOutPut("设置成功", "设置成功!")
	return
}

func MailInfo(c *ginHelper.GinCtx) {
	data, err := model.GetMail()
	if err != nil {
		c.APIOutPutError(nil, err.Error())
		return
	}
	c.APIOutPut(data, "")
	return
}

type MailSendTestParam struct {
	To    string `json:"to"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func MailSendTest(c *ginHelper.GinCtx) {
	err := mailSet(c)
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	title := "站点监控发送邮件通知测试"
	body := `<p>您好欢迎使用小型站点监测平台(Small website monitor)</p>` +
		`<p> 您的星星是我研发的动力!</p>` +
		`<p><a herf="https://github.com/mangenotwork/small-website-monitor">https://github.com/mangenotwork/small-website-monitor</a></p>` +
		`<p>------ ManGe ` + time.Now().String() + `</p>`
	model.Send(title, body)
	c.APIOutPut("", "测试邮件已发送请注意查收!")
	return
}

type WebsitePointParam struct {
	Uri string `json:"uri"`
}

func WebsitePointAdd(c *ginHelper.GinCtx) {
	hostId := c.Param("hostId")
	log.Info("hostId = ", hostId)
	param := &WebsitePointParam{}
	err := c.GetPostArgs(param)
	if err != nil {
		c.APIOutPutError(err, "参数或参数类型错误")
		return
	}
	ctx, _ := gt.Get(param.Uri)
	if business.AlertRuleCode(ctx.StateCode) {
		c.APIOutPutError(nil, fmt.Sprintf("%s请求失败，状态码:%d", param.Uri, ctx.StateCode))
		return
	}
	websitePoint := model.NewWebSitePoint(hostId)
	err = websitePoint.Add(param.Uri)
	if err != nil {
		c.APIOutPutError(err, "添加监测点失败:"+err.Error())
		return
	}
	c.APIOutPut("", "添加成功")
	return
}

func WebsitePointList(c *ginHelper.GinCtx) {
	hostId := c.Param("hostId")
	websitePoint := model.NewWebSitePoint(hostId)
	err := websitePoint.Get()
	if err != nil {
		c.APIOutPutError(err, "获取失败")
		return
	}
	c.APIOutPut(websitePoint.Uri, "")
	return
}

func WebsitePointDel(c *ginHelper.GinCtx) {
	hostId := c.Param("hostId")
	param := &WebsitePointParam{}
	err := c.GetPostArgs(param)
	if err != nil {
		c.APIOutPutError(err, "参数或参数类型错误")
		return
	}
	websitePoint := model.NewWebSitePoint(hostId)
	err = websitePoint.Del(param.Uri)
	if err != nil {
		c.APIOutPutError(err, "删除失败:"+err.Error())
		return
	}
	c.APIOutPut(websitePoint.Uri, "删除成功")
	return
}

type WebsiteInfoOut struct {
	*model.WebSite
	*model.WebSiteUri
}

func WebsiteInfo(c *ginHelper.GinCtx) {
	hostId := c.Param("hostId")
	website, err := new(model.WebSite).Get(hostId)
	if err != nil {
		c.APIOutPutError(err, "获取失败:"+err.Error())
		return
	}
	websiteList := model.NewWebSiteUri(hostId)
	_, _ = websiteList.Get()
	data := &WebsiteInfoOut{website, websiteList}
	c.APIOutPut(data, "")
	return
}

func AlertList(c *ginHelper.GinCtx) {
	websiteList, err := new(model.WebSite).GetAll()
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	data := make([]string, 0)
	for _, v := range websiteList {
		alert := model.NewWebSiteAlert(v.ID)
		err = alert.Get()
		if err != nil {
			log.Error(err)
			continue
		}
		for _, a := range alert.List {
			data = append(data, fmt.Sprintf("%s : %s | %s", a.Date, a.Uri, a.Msg))
		}
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i] > data[j]
	})
	c.APIOutPut(data, "")
	return
}

func AlertClear(c *ginHelper.GinCtx) {
	websiteList, err := new(model.WebSite).GetAll()
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	var hasErr error
	for _, v := range websiteList {
		alert := model.NewWebSiteAlert(v.ID)
		hasErr = alert.Clear()
	}
	if err != nil {
		c.APIOutPutError(hasErr, "清空失败，err = "+hasErr.Error())
		return
	}
	c.APIOutPut("", "清空完成")
	return
}

func MonitorErrList(c *ginHelper.GinCtx) {
	data := model.NewMonitorErrInfo()
	err := data.Get()
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	sort.Slice(data.List, func(i, j int) bool {
		return data.List[i] > data.List[j]
	})
	c.APIOutPut(data.List, "")
	return
}

func MonitorErrClear(c *ginHelper.GinCtx) {
	data := model.NewMonitorErrInfo()
	err := data.Clear()
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	c.APIOutPut("", "清空成功")
	return
}

func MonitorLog(c *ginHelper.GinCtx) {
	hostId := c.Param("hostId")
	data := business.ReadLog(hostId)
	c.APIOutPut(data, "")
	return
}

func WebsiteDelete(c *ginHelper.GinCtx) {
	hostId := c.Param("hostId")
	// 确认是否存在
	website, err := new(model.WebSite).Get(hostId)
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	if utils.AnyToInt64(website.ID) < 1 {
		c.APIOutPutError(fmt.Errorf("站点不存在"), "站点不存在")
		return
	}
	// 删除website
	err = website.Delete(website.ID)
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	// 删除Uri
	websiteUri := model.NewWebSiteUri(hostId)
	err = websiteUri.Delete()
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	// 删除Point
	websitePoint := model.NewWebSitePoint(hostId)
	err = websitePoint.DeleteWebsite()
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	// 更新website对象
	business.Push()
	// 删除日志
	err = business.DeleteLog(website.ID)
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	c.APIOutPut("", "删除成功")
	return
}

func AlertCount(c *ginHelper.GinCtx) {
	hostId := c.Param("hostId")
	alert := model.NewWebSiteAlert(hostId)
	alertCount, err := alert.Len()
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	c.APIOutPut(alertCount, "")
	return
}

type WebsiteEditParam struct {
	HostId       string `json:"hostId"`
	Rate         int64  `json:"rate"`
	AlarmResTime int64  `json:"alarmResTime"`
	UriDepth     int64  `json:"uriDepth"`
}

func WebsiteEdit(c *ginHelper.GinCtx) {
	param := &WebsiteEditParam{}
	err := c.GetPostArgs(param)
	if err != nil {
		c.APIOutPutError(err, "参数或参数类型错误")
		return
	}
	if param.Rate < 1 {
		c.APIOutPutError(nil, "监测频率不能小于1秒")
		return
	}
	if param.UriDepth < 1 {
		c.APIOutPutError(nil, "采集URI深度不能小于1层")
		return
	}
	if param.AlarmResTime < 1 {
		c.APIOutPutError(nil, "报警时间不能小于1ms")
		return
	}
	// 更新设置
	website, err := new(model.WebSite).Get(param.HostId)
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	website.Rate = param.Rate
	website.UriDepth = param.UriDepth
	website.AlarmResTime = param.AlarmResTime
	err = website.Update()
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	// 重置对象
	business.Push()
	c.APIOutPut("", "设置成功")
	return
}

func WebsiteChart(c *ginHelper.GinCtx) {
	hostId := c.Param("hostId")
	day := c.GetQuery("day")
	if len(day) < 1 {
		day = utils.NowDateLayout(global.DayLayout)
	}
	uri := c.GetQuery("uri")
	data, err := business.ReadAll(hostId, day)
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	out := make([][]int64, 0)
	for _, v := range data {
		if v.UriType == uri {
			item := make([]int64, 0)
			item = append(item, utils.Date2Timestamp(v.Time)*1000)
			item = append(item, v.UriMs)
			out = append(out, item)
		}
	}
	c.APIOutPut(out, "")
	return
}

func WebsiteAlertList(c *ginHelper.GinCtx) {
	hostId := c.Param("hostId")
	alert := model.NewWebSiteAlert(hostId)
	err := alert.Get()
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	c.APIOutPut(alert.List, "")
	return
}

func WebsiteAlertDel(c *ginHelper.GinCtx) {
	hostId := c.Param("hostId")
	date := c.GetQuery("date")
	alert := model.NewWebSiteAlert(hostId)
	err := alert.Get()
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	log.Info("date = ", date)
	for n, v := range alert.List {
		if v.Date == date {
			log.Info("删除.....")
			alert.List = append(alert.List[:n], alert.List[n+1:]...)
			break
		}
	}
	err = alert.Update()
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	c.APIOutPut(hostId, "成功")
	return
}

type WebsiteLogListOut struct {
	FileList []string
	DayList  []string
}

func WebsiteLogList(c *ginHelper.GinCtx) {
	hostId := c.Param("hostId")
	data, err := business.LogList(hostId)
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i] > data[j]
	})
	day := make([]string, 0)
	for _, v := range data {
		d := strings.Split(v, "_")
		if len(d) == 2 {
			day = append(day, strings.Replace(d[1], ".log", "", -1))
		}
	}
	out := &WebsiteLogListOut{
		FileList: data,
		DayList:  day,
	}
	c.APIOutPut(out, "")
	return
}

func WebsiteLogUpload(c *ginHelper.GinCtx) {
	hostId := c.Param("hostId")
	day := c.GetQuery("day")
	logPath, err := business.Upload(hostId, day)
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	c.Writer.Header().Add("Content-Disposition",
		fmt.Sprintf("attachment; filename=%s", fmt.Sprintf("%s.log", day)))
	c.Writer.Header().Add("Content-Type", "application/text/plain")
	c.File(logPath)
	return
}

func Case1(c *ginHelper.GinCtx) {
	id, err := model.GetIncrement()
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	c.APIOutPut(id, "")
	return
}

func Case2(c *ginHelper.GinCtx) {
	err := model.ResetIncrement()
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	c.APIOutPut("", "ok")
	return
}
