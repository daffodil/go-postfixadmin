package postfixadmin

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/daffodil/go-postfixadmin/base"
)

type VacationsPayload struct {
	Success   bool        ` json:"success" `
	Vacations []*Vacation ` json:"vacations" `
	Error     string      ` json:"error" `
}

func GetVacations(domain string) ([]*Vacation, error) {
	var rows []*Vacation
	var err error

	if DomainExists(domain) == false {
		return rows, errors.New("Domain `" + domain + "` does not exist")
	}

	Dbo.Where("domain=?", domain).Find(&rows)
	return rows, err
}

// /domain/<domain>/vacations
func HandleAjaxVacations(resp http.ResponseWriter, req *http.Request) {

	if base.AjaxAuth(resp, req) == false {
		return
	}

	fmt.Println("VacationsAjaxHandler")
	vars := mux.Vars(req)

	payload := VacationsPayload{}
	payload.Success = true
	payload.Vacations = make([]*Vacation, 0)

	var err error
	payload.Vacations, err = GetVacations(vars["domain"])
	if err != nil {
		fmt.Println(err)
		payload.Error = "DB Error: " + err.Error()
	}

	base.SendPayload(resp, payload)

}
