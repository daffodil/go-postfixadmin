

package sendmail

import (

	"fmt"
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

func SendMessage(mess Message) error {
	// Set up authentication information.
	fmt.Println("SendMessage()", mess.Subject)
	m := gomail.NewMessage()

	m.SetHeader("From", mess.From)
	m.SetHeader("To", mess.To...)
	m.SetHeader("Bcc", Conf.SyslogEmail)

	m.SetHeader("Subject", mess.Subject)
	m.SetBody("text/plain", mess.Body)

	c := gomail.NewPlainDialer(Conf.SMTPServer.Server, Conf.SMTPServer.Port, Conf.SMTPServer.Login, Conf.SMTPServer.Password)
	c.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err := c.DialAndSend(m)
	if err != nil {
		return err
	}
	return nil
}