

package sendmail


import (

	"fmt"



)

func SendAdminMessage(subject string, body string) error {
	// Set up authentication information.
	fmt.Println("SendAdminMessage()", subject)

	m :=  Message{}

	m.From = Conf.FromEmail
	m.AddTo( Conf.AdminEmail)
	m.AddBcc( Conf.SyslogEmail)
	//m.SetHeader("Bcc", Conf.SyslogEmail)

	m.SetSubject( Conf.EmailPrefix + subject)
	m.SetBody("text/plain", body)


	return SendMessage(m)
	/*
	c := gomail.NewPlainDialer(Conf.Server, Conf.Port, Conf.Login, Conf.Password)
	c.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err := c.DialAndSend(m)
	if err != nil {
		return err
	}
	*/
	//return nil
}