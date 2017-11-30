
package sendmail

import (


)


type Message struct {
	From string
	To []string
	Cc []string
	Bcc []string
	Subject string
	ContentType string
	Body string
}


func (m *Message) AddTo(addr string) {
	m.To = append(m.To, addr)
}

func (m *Message) AddCc(addr string) {
	m.Cc = append(m.Cc, addr)
}

func (m *Message) AddBcc(addr string) {
	m.Bcc = append(m.Bcc, addr)
}
func (m *Message) SetSubject(subject string) {
	m.Subject = subject

}
func (m *Message) SetBody(content_type string, content string) {
	m.ContentType = content_type
	m.Body = content
}
