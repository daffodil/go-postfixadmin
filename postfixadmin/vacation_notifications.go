package postfixadmin

import (

	"net/http"

	"github.com/daffodil/go-postfixadmin/base"
)

type VacationNotification struct {
	OnVacation string ` json:"on_vacation" `
	Notified   string ` json:"notified" `
	NotifiedAt string ` json:"notified_at" `
}

func GetVacationNotifications(email string, order string) ([]VacationNotification, error) {
	var rows []VacationNotification
	var err error
	var order_by string
	if order == "email" {
		order_by = "notified asc"
	} else {
		order_by = "notified_at desc"
	}
	Dbo.Where("on_vacation = ?", email).Order(order_by).Limit(100).Find(&rows)
	return rows, err
}

func GetRecentVacationNotifications() ([]VacationNotification, error) {
	var rows []VacationNotification
	var err error
	Dbo.Order("notified_at desc").Limit(100).Find(&rows)
	return rows, err
}

type VacationNotificationsPayload struct {
	Success       bool                   `json:"success"` // keep extjs happy
	Notifications []VacationNotification `json:"vacation_notifications"`
	Error         string                 `json:"error"`
}

// Handles /ajax/vacation/<email>
func HandleAjaxVacationNotifications(resp http.ResponseWriter, req *http.Request) {

	if base.AjaxAuth(resp, req) == false {
		return
	}

	payload := VacationNotificationsPayload{}
	payload.Success = true

	var err error
	payload.Notifications, err = GetRecentVacationNotifications()
	if err != nil {
		payload.Error = err.Error()
	}

	base.SendPayload(resp, payload)
}
