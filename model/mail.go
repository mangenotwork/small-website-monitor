package model

import (
	"fmt"
	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/mail"
)

type MailData struct {
	From     string
	AuthCode string
	Host     string
	Port     int
	ToList   []string
}

func IsMail() bool {
	value := &MailData{}
	err := DB.Get(MailTable, MailConf, value)
	if err != nil {
		log.Error(err)
	}
	log.Info("value = ", value)
	if len(value.Host) > 0 && len(value.From) > 0 && len(value.AuthCode) > 0 {
		return true
	}
	return false
}

func SetMail(mailConf *MailData) error {
	return DB.Set(MailTable, MailConf, mailConf)
}

func GetMail() (*MailData, error) {
	value := &MailData{}
	err := DB.Get(MailTable, MailConf, value)
	return value, err
}

func Send(title, body string) {
	mailObj, err := GetMail()
	if err != nil {
		log.Error(err)
		return
	}
	m := mail.NewMail(mailObj.Host, mailObj.From, mailObj.AuthCode, mailObj.Port)
	err = m.Title(title).HtmlBody(body).SendMore(mailObj.ToList)
	if err != nil {
		log.Error(err)
		return
	}
}

func SendTest(title, body string, to []string) {
	mailObj, err := GetMail()
	if err != nil {
		log.Error(err)
		return
	}
	m := mail.NewMail(mailObj.Host, mailObj.From, mailObj.AuthCode, mailObj.Port)
	err = m.Title(title).HtmlBody(body).SendMore(to)
	if err != nil {
		log.Error(err)
		return
	}
}

// AlertBody 报警通知
type AlertBody struct {
	Synopsis string
	Tr       []*AlertTd
}

type AlertTd struct {
	Date string
	Host string
	Uri  string
	Code int
	Ms   int64
	Msg  string
}

func (a *AlertBody) Html() string {
	body := ""
	synopsis := fmt.Sprintf("<h3>%s</h3>", a.Synopsis)
	tr := ""
	for _, v := range a.Tr {
		tr += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%d</td><td>%d</td><td>%sms</td></tr>",
			v.Date, v.Host, v.Uri, v.Code, v.Ms, v.Msg)
	}
	thead := "<thead><tr><td>监测时间</td><td>站点</td><td>链接</td><td>请求状态码</td><td>响应时间</td><td>报警信息</td></tr></thead>"
	table := fmt.Sprintf("<table>%s<tbody>%s</tbody></table>", thead, tr)
	body = synopsis + table
	return body
}
