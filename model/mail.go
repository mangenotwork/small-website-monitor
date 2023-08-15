package model

import (
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
