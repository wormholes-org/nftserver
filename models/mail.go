package models

import (
	"crypto/tls"
	"fmt"
	"github.com/go-gomail/gomail"
	"log"
)

type Cc struct {
	address string
	name    string
}

type MailInfo struct {
	From    string
	To      []string
	Carbon  []Cc
	Subject string
	Body    string
}

type MailSvr struct {
	Host     string
	Port     int
	UserName string
	PassWord string
}

const (
	EmailBodyModel    = "Hello, the exchange you handled in %v will expire on %v please log in to the website in time to renew online, thank you for your support to Wormholes."
	EmailSubjectModel = "Service expires"
	EmailHost         = "mail.wormholes.com"
	EmailPort         = 25
	EmailUserName     = "li.haisheng"
	EmailUserPassword = "Wormholes_&639"
)

func SendMail(mailinfo *MailInfo, mailsvr *MailSvr) error {
	m := gomail.NewMessage()
	m.SetHeader("From", mailinfo.From)
	m.SetHeader("To", mailinfo.To...)
	var ccs []string
	for _, cc := range mailinfo.Carbon {
		ccs = append(ccs, m.FormatAddress(cc.address, cc.name))
	}
	m.SetHeader("Cc", ccs...)
	m.SetHeader("Subject", mailinfo.Subject)
	m.SetBody("text/html", mailinfo.Body)
	//m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer(mailsvr.Host, mailsvr.Port, mailsvr.UserName, mailsvr.PassWord)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		log.Println("SendMail() DialAndSend err=", err)
		return err
	}
	return nil
}

func (nft *NftDb) TranMail(useraddr string) error {
	var user Users
	db := nft.db.Model(&Users{}).Where("useraddr= ?", useraddr).First(&user)
	if db.Error != nil {
		log.Println("TranMail first err =", db.Error)
		return db.Error
	}
	mailInfo := MailInfo{}
	mailInfo.From = useraddr
	mailInfo.Subject = EmailSubjectModel
	//mailInfo.Body = agent.Emailbody
	mailInfo.To = []string{}
	//mailInfo.Carbon = []Cc{{agent.Adminemail, agent.Emailaccount}}
	mailSvr := MailSvr{}
	mailSvr.Host = EmailHost
	mailSvr.Port = EmailPort
	mailSvr.UserName = EmailUserName
	mailSvr.PassWord = EmailUserPassword
	mailInfo.To = append(mailInfo.To, user.Email)
	emailbody := fmt.Sprintf(EmailBodyModel, "", "")
	mailInfo.Body = emailbody
	err := SendMail(&mailInfo, &mailSvr)
	if err != nil {
		log.Println("TranMail SendMail err=", err)
		return err
	}
	return nil
}
