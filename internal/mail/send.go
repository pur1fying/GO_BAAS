package mail

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"

	"github.com/pur1fying/GO_BAAS/internal/config"
	"github.com/pur1fying/GO_BAAS/internal/logger"
)

var tlsConfig = &tls.Config{
	InsecureSkipVerify: true,
}

var auth smtp.Auth

var _SMTPClient *smtp.Client

var mailQueue chan *Message

var smtpAddr string

var smtpPort int

var smtpHost string

var userName string

var authCode string

func InitMail() error {
	smtpPort = config.Config.Mail.SMTPPort
	smtpHost = config.Config.Mail.SMTPHost
	userName = config.Config.Mail.Username
	authCode = config.Config.Mail.AuthCode
	smtpAddr = fmt.Sprintf("%s:%d", smtpHost, smtpPort)

	logger.HighLight("Init Mail")
	logger.BAASInfo("Addr :", smtpAddr)
	if err := checkPort(); err != nil {
		return err
	}

	return DialAndAuth()
}

func checkPort() error {
	if smtpPort == 465 || smtpPort == 476 {
		return nil
	}
	return errors.New("Invalid port " + fmt.Sprint(smtpPort) + " for SMTP. Only 465 and 587 are supported")

}

func InitSMTPAuth() error {
	auth = smtp.PlainAuth("", userName, authCode, smtpHost)
	return _SMTPClient.Auth(auth)
}

func Dial() (err error) {
	var conn net.Conn
	if smtpPort == 465 {
		conn, err = tls.Dial("tcp", smtpAddr, nil)
	} else if smtpPort == 587 {
		conn, err = net.Dial("tcp", smtpAddr)
	}

	if err != nil {
		logger.BAASError("Dial SMTP Error : ", err.Error())
		return err
	}

	_SMTPClient, err = smtp.NewClient(conn, smtpHost)
	if err != nil {
		return errors.New("Create SMTP Client Error : " + err.Error())
	}

	if smtpPort == 587 {
		err = _SMTPClient.Hello("localhost")
		if err != nil {
			return errors.New("HELO Failed : " + err.Error())
		}

		if ok, _ := _SMTPClient.Extension("STARTTLS"); ok {
			err = _SMTPClient.StartTLS(tlsConfig)
			if err != nil {
				return errors.New("STARTTLS Failed : " + err.Error())
			}
		} else {
			logger.BAASWarn(smtpHost, "Do Not Support STARTTLS")
		}
	}

	return nil
}

func SendMail(msg *Message) error {
	// Mail
	if err := _SMTPClient.Mail(msg.From.Email); err != nil {
		return err
	}

	// Rcpt
	for _, addr := range msg.TO {
		if err := _SMTPClient.Rcpt(addr.Email); err != nil {
			return err
		}
	}

	for _, addr := range msg.CC {
		if err := _SMTPClient.Rcpt(addr.Email); err != nil {
			return err
		}
	}

	for _, addr := range msg.BCC {
		if err := _SMTPClient.Rcpt(addr); err != nil {
			return err
		}
	}

	// Data
	w, err := _SMTPClient.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(msg.generateHeader()))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(msg.generateBody()))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}

func DialAndAuth() error {
	if err := Dial(); err != nil {
		return err
	}
	logger.BAASInfo("SMTP Dial Success")
	if err := InitSMTPAuth(); err != nil {
		return err
	}
	logger.BAASInfo("SMTP Auth Success")
	return nil
}
