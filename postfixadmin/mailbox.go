package postfixadmin

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/daffodil/go-postfixadmin/base"
	"github.com/daffodil/go-postfixadmin/sendmail"
)

type Mailbox struct {
	Username string `json:"username" gorm:"primary_key"`
	//Password string 	`json:"password"`
	Name      string `json:"name"`
	Maildir   string `json:"maildir"`
	Quota     int    `json:"quota"`
	LocalPart string `json:"local_part"`
	Domain    string `json:"domain"`
	Created   string `json:"created"`
	Modified  string `json:"modified"`
	Active    bool   `json:"active"`
}

// gorm request table name here
func (me Mailbox) TableName() string {
	return Conf.Db.TableNames["mailbox"]
}

// Updates a mailboxes password
func SetMailboxPassword(email string, secret string) error {

	if MailboxExists(email) == false {
		//TODO bill make this into official error predefined above
		return errors.New("Mailbox not exist")
	}
	sql := "update mailbox set password=? where username = ?"
	Dbo.Exec(sql, secret, email)
	return nil
}

// Retrieves the unencryped pass from DB , cos its a pain i the ass with tokens, so raw for now..
// same caper i postfix, so we expect behind fire wall
func GetMailboxPassword(email string) (string, error) {

	if MailboxExists(email) == false {
		//TODO bil make error
		return "", errors.New("Mailbox not exist")
	}
	var pass string
	sql := "select password from mailbox where username = ?"
	row := Dbo.Raw(sql, email).Row()
	row.Scan(&pass)
	return pass, nil
}

// Load a mailBox record fro the ORM database
func GetMailbox(username string) (Mailbox, error) {
	var mailbox Mailbox
	var err error
	Dbo.Where("username = ? ", username).First(&mailbox)
	return mailbox, err
}

// Checks a mailbox exists
func MailboxExists(address string) bool {
	var count int
	Dbo.Model(Mailbox{}).Where("username = ?", address).Count(&count)
	if count == 0 {
		return false
	}
	return true
}

// Struct for enconding JAon and sending to client
type MailboxPayload struct {
	Success bool    `json:"success"` // keep extjs happy
	Mailbox Mailbox `json:"mailbox"`
	Aliases string  `json:"aliases"`
	Error   string  `json:"error"`
}

// Struct container for sending payload..
// TODO bill make tis a creeatePayload an map ?? ummm .. arg
// this is specific and mapped and fast  maybe
func CreateMailboxPayload() MailboxPayload {
	payload := MailboxPayload{}
	payload.Success = true
	payload.Mailbox = Mailbox{}
	return payload
}

// /domain/<example.com>/mailbox/<email>
func HandleAjaxMailbox(resp http.ResponseWriter, req *http.Request) {

	//create ajax payload
	payload := CreateMailboxPayload()

	// grab vars from gorilla and in the URL
	vars := mux.Vars(req)
	domain := vars["domain"]
	username := vars["username"]

	// check domain valid
	errdom := IsDomainValid(domain)
	if errdom != nil {
		payload.Error = errdom.Error()
		base.WriteJSON(resp, payload)
		return
	}

	// check email address valid
	email_addr, err_email := ParseAddress(username)
	if err_email != nil {
		payload.Error = err_email.Error()
		base.WriteJSON(resp, payload)
		return
	}

	// check email and domain match
	if payload.Error != "" && email_addr.Domain != domain {
		payload.Error = errors.New(" Domain and email domain do not match").Error()
		base.WriteJSON(resp, payload)
		return
	}

	// Check mailbox etc
	var err error
	exists := MailboxExists(email_addr.Address)
	if exists {
		payload.Mailbox, err = GetMailbox(email_addr.Address)
	}
	if err != nil {
		payload.Error = "DB Error: " + err.Error()
		base.WriteJSON(resp, payload)
		return
	}

	var alias *Alias

	switch req.Method {

	case "POST":
		// TODO: REST , currently is uses POST and form vars
		form := req.Form

		payload.Mailbox.Username = email_addr.Address
		payload.Mailbox.Domain = email_addr.Domain
		payload.Mailbox.LocalPart = email_addr.User
		payload.Mailbox.Active, _ = strconv.ParseBool(form.Get("active"))
		payload.Mailbox.Name = form.Get("name")
		payload.Mailbox.Maildir = domain + "/" + email_addr.User
		payload.Mailbox.Quota = 0
		if payload.Mailbox.Created == "" {
			payload.Mailbox.Created = base.NowStr()
		}
		payload.Mailbox.Modified = base.NowStr()

		// If mailbox not exist we create
		if exists == false {
			Dbo.Create(&payload.Mailbox)
		} else {
			Dbo.Save(&payload.Mailbox)
		}

		// This is a mad quirk for postfixadmin..
		// But all emails are an alias (from high performance systems apparently)
		if AliasExists(email_addr.Address) == false {
			alias = &Alias{}
			alias.Address = email_addr.Address
			alias.Domain = email_addr.Domain
			alias.Created = base.NowStr()
			Dbo.Create(&alias)
		} else {
			alias, _ = GetAlias(email_addr.Address)
		}

		alias.ClearAllGoto()
		aliases_raw := form.Get("aliases")
		aliases := strings.Split(aliases_raw, ",")
		fmt.Println("####", aliases_raw, aliases)
		for _, a := range aliases {
			alias.AddGoto(a)
		}
		alias.Save()

		// WOW dodgy////
		if form.Get("password") != "" {
			SetMailboxPassword(payload.Mailbox.Username, form.Get("password"))
		}

		payload.Aliases = alias.Goto

	case "GET":

		alias, _ = GetAlias(email_addr.Address)
		payload.Aliases = alias.Goto
		// just pass through vars

	}

	base.WriteJSON(resp, payload)
}

type MailboxPassPayload struct {
	Success bool   `json:"success"` // keep extjs happy
	Ok      bool   `json:"OK"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

func HandleAjaxSetMailboxPassword(resp http.ResponseWriter, req *http.Request) {


	payload := MailboxPassPayload{}
	payload.Success = true

	//vars := mux.Vars(req)
	f := req.Form

	email_raw := f.Get("email")
	email_addr, err_email := ParseAddress(email_raw)
	if err_email != nil {
		payload.Error = err_email.Error()
		payload.Ok = false
		base.WriteJSON(resp, payload)
		return
	}

	passwd := f.Get("secret")
	if len(passwd) < 4 {
		payload.Error = "Password too short"
	} else {

		SetMailboxPassword(email_addr.Address, passwd)
		payload.Ok = true
		payload.Message = "Password Set"
	}

	base.WriteJSON(resp, payload)

}

type SendTestMailPayload struct {
	Success bool   `json:"success"` // keep extjs happy
	Ok      bool   `json:"OK"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

func HandleAjaxMailboxSendTest(resp http.ResponseWriter, req *http.Request) {


	payload := SendTestMailPayload{}
	payload.Success = true

	vars := mux.Vars(req)
	email_raw := vars["email"]
	email_addr, err_email := ParseAddress(email_raw)
	if err_email != nil {
		payload.Error = err_email.Error()
		payload.Ok = false
		base.WriteJSON(resp, payload)
		return
	}

	tim := time.Now()
	t := tim.String()

	mess := sendmail.Message{}
	mess.AddTo(email_addr.Address)
	mess.From = Conf.NoReplyEmail
	mess.Subject = "Test Message: " + t
	mess.Body = "Test" + t

	err := sendmail.SendMessage(mess)
	if err != nil {
		payload.Error = err.Error()
	}

	base.WriteJSON(resp, payload)
}
