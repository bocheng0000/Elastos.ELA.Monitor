package email

import (
	"bytes"
	"fmt"
	"github.com/elastos/Elastos.ELA.Monitor/config"
	"github.com/elastos/Elastos.ELA.Monitor/utility/log"
	"net/smtp"
	"strings"
)

func SendMonitorEMail(subject, content string) error {
	var err error
	emailConfig := &config.ConfigManager.MonitorConfig.EMail
	for _, notifyUser := range (*emailConfig).NotifyUser {
		err = SendEMail("", (*emailConfig).Host, (*emailConfig).UserName, (*emailConfig).PassWord, notifyUser, subject, content)
		if err != nil {
			log.Errorf("send email to %s failed", notifyUser)
		}
	}

	return err
}

func SendEMail(identity, host, userName, passWord, toUser, subject, content string) error {
	var messageBuffer bytes.Buffer
	messageBuffer.WriteString(fmt.Sprintf("From: %s\r\n", userName))
	messageBuffer.WriteString(fmt.Sprintf("To: %s\r\n", toUser))
	messageBuffer.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	messageBuffer.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
	messageBuffer.WriteString(content)

	plainAuth := smtp.PlainAuth(identity, userName, passWord, strings.Split(host, ":")[0])
	err := smtp.SendMail(host, plainAuth, userName, []string{ toUser }, []byte(messageBuffer.String()))
	if err == nil {
		log.Infof("send email to %s successed", toUser)
	}

	return err
}
