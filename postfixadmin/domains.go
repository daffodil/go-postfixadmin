package postfixadmin

import (
	"fmt"
	"net/http"

	"github.com/daffodil/go-postfixadmin/base"
)

// Load domains from database
func LoadDomains() ([]Domain, error) {
	var rows []Domain
	var err error
	Dbo.Where("domain <> ?", "ALL").Find(&rows)
	return rows, err
}

// The ajax struct send as json
type DomainsPayload struct {
	Success bool     `json:"success"` // keep extjs happy
	Domains []Domain `json:"domains"`
	Error   string   `json:"error"`
}

// /domains
func HandleAjaxDomains(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("DomainsAjaxHandler")
	if base.AjaxAuth(resp, req) == false {
		return
	}

	payload := DomainsPayload{}
	payload.Success = true
	//t.Domains = make([]Domain, 0)

	var err error
	payload.Domains, err = LoadDomains()
	if err != nil {
		fmt.Println(err)
		payload.Error = "DB Error: " + err.Error()
	}

	base.SendPayload(resp, payload)
}