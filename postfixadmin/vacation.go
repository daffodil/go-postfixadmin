package postfixadmin

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/daffodil/go-postfixadmin/base"
)

type Vacation struct {
	Email        string ` json:"email" gorm:"primary_key" `
	Subject      string ` json:"subject" `
	Body         string ` json:"body" `
	Activefrom   string ` json:"active_from" `
	Activeuntil  string ` json:"active_until" `
	Cache        string ` json:"cache" `
	Domain       string ` json:"domain" `
	IntervalTime int64  ` json:"interval_time" `
	Created      string ` json:"created" `
	Modified     string ` json:"modified" `
	Active       bool   ` json:"active" `
	LastBot      string ` json:"last_bot" `
}

// Return payload for vacation
type VacationPayload struct {
	Success               bool                   `json:"success"` // keep extjs happy
	Vacation              *Vacation              `json:"vacation"`
	VacationNotifications []VacationNotification ` json:"vacation_notifications" `
	Error                 string                 `json:"error"`
}

// return true is address is a vacation address
func IsVacationAddress(address string) bool {
	if address == "" {
		return false
	}
	user_domain := strings.Split(address, "@")
	fmt.Println("-------------------", address, user_domain)
	if len(user_domain) == 0 {
		return false
	}
	if user_domain[1] == Conf.VacationDomain {
		return true
	}
	return false
}

// Load vacation record from db
func GetVacation(email string) (*Vacation, error) {
	row := new(Vacation)
	var err error
	Dbo.Where("email = ?", email).Find(&row)
	return row, err
}

// Check is an email exists in the vacation table
func VacationExists(address string) bool {
	var count int
	Dbo.Model(Vacation{}).Where("email = ?", address).Count(&count)
	if count == 0 {
		return false
	}
	return true
}

// /vacation/<email>
func HandleAjaxVacation(resp http.ResponseWriter, req *http.Request) {

	payload := VacationPayload{}
	payload.Success = true //extjs fu

	vars := mux.Vars(req)

	email_addr, err_email := ParseAddress(vars["email"])
	if err_email != nil {
		payload.Error = err_email.Error()
	} else {

		// check mail exists
		if !MailboxExists(email_addr.Address) {
			payload.Error = errors.New("Mailbox `" + email_addr.Address + "` does not exist").Error()

		} else {

			var err error
			if VacationExists(email_addr.Address) {

				payload.Vacation, err = GetVacation(email_addr.Address)
				if err != nil {
					fmt.Println(err)
					payload.Error = "DB Error: " + err.Error()
				}
			} else {
				payload.Vacation = new(Vacation)
				payload.Vacation.Email = email_addr.Address
				payload.Vacation.Domain = email_addr.Domain
				Dbo.Create(&payload.Vacation)
			}
			switch req.Method {

			case "POST":

				f := req.Form
				fmt.Println(f)

				payload.Vacation.Active, err = strconv.ParseBool(f.Get("active"))
				payload.Vacation.Activefrom = f.Get("active_from")
				payload.Vacation.Activeuntil = f.Get("active_until")
				payload.Vacation.IntervalTime, err = strconv.ParseInt(f.Get("interval_time"), 10, 64)
				payload.Vacation.Subject = f.Get("subject")
				payload.Vacation.Body = f.Get("body")

				Dbo.Save(&payload.Vacation)
				fmt.Println("------------------ POSTED-------------", f.Get("active_from"))
				UpdateVacationAlias(payload.Vacation)

			case "GET":

				if payload.Vacation.Email == "" {
					// probably record not exist
					payload.Vacation.Email = email_addr.Address
				} else {
					payload.VacationNotifications, err = GetVacationNotifications(email_addr.Address, "date")
				}
			}
		}
	}

	base.WriteJSON(resp, payload)
}

// Updates a Vacation record
func UpdateVacationAlias(vac *Vacation) {

	alias, err := GetAlias(vac.Email)
	fmt.Println("UpdateVacationAlias", alias, err)
	if err != nil {
		// do something
		return
	}
	em, errp := ParseAddress(vac.Email)
	if errp != nil {
		return
	}
	if vac.Active {
		alias.AddGoto(em.VacationAddress)
	} else {
		alias.RemoveGoto(em.VacationAddress)
	}
	alias.Save()
}
