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
			c.HTML(200, "home.html", gin.H{})
			return
		}
	}
	c.HTML(200, "err.html", gin.H{
		"err": "账号或密码错误!",
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
	MoreUri        string `json:"moreUri"`
}

func WebsiteAdd(c *ginHelper.GinCtx) {
	param := &WebsiteAddParam{}
	err := c.GetPostArgs(param)
	if err != nil {
		c.APIOutPutError(err, err.Error())
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
		HealthUri:     param.HealthUri,
		Rate:          param.Rate,
		UriDepth:      param.UriDepth,
		UriUpdateRate: param.UriUpdateRate,
		AlarmResTime:  param.AlarmResTime,
		HostIP:        ctx.Req.RemoteAddr,
		Created:       time.Now().Unix(),
	}
	log.Info(website)
	err = website.Add()
	if err != nil {
		c.APIOutPutError(err, "保存数据失败")
		return
	}

	// 更新站点对象
	business.Push()

	c.APIOutPut("", "添加成功")
	return
}

func CaseT(c *ginHelper.GinCtx) {
	c.APIOutPut("test", "")
}

type WebsiteListOut struct {
	List     []*model.WebSite
	Count    int
	PageList []*ginHelper.Page
}

func WebsiteList(c *ginHelper.GinCtx) {

	pg := c.GetQueryInt("pg")
	if pg < 1 {
		pg = 1
	}
	log.Info("pg = ", pg)
	size := 10

	// TODO 聚合查询，需要将最近响应时间取出来聚合

	data, count, err := new(model.WebSite).List(pg, size)
	if err != nil {
		c.APIOutPutError(err, "获取失败")
		return
	}

	c.APIOutPut(&WebsiteListOut{
		List:     data,
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

func MailConf(c *ginHelper.GinCtx) {
	param := &MailConfParam{}
	err := c.GetPostArgs(param)
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	if len(param.From) < 1 {
		c.APIOutPutError(nil, "发件人不能为空!")
		return
	}
	if len(param.Host) < 1 {
		c.APIOutPutError(nil, "邮件服务器不能为空!")
		return
	}
	if len(param.AuthCode) < 1 {
		c.APIOutPutError(nil, "邮件服务授权码不能为空!")
		return
	}
	if len(param.ToList) < 1 {
		c.APIOutPutError(nil, "通知收件人不能为空!")
		return
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
	err = model.SetMail(m)
	if err != nil {
		c.APIOutPutError(nil, err.Error())
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
	title := "站点监控发送邮件通知测试"
	body := `<p>您好欢迎使用小型站点监测平台(Small website monitor)</p>` +
		`<p> 您的星星是我研发的动力!</p>` +
		`<p><a herf="https://github.com/mangenotwork/small-website-monitor">https://github.com/mangenotwork/small-website-monitor</a></p>` +
		`<p>------ ManGe ` + time.Now().String() + `</p>`
	model.Send(title, body)
	c.APIOutPut("", "测试邮件已发送请注意查收!")
	return
}
