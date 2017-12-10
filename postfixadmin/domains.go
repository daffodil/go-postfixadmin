package postfixadmin

import (
	"fmt"
	"net/http"

	"github.com/daffodil/go-postfixadmin/base"
)

// Load domains from database
func GetDomains() ([]Domain, error) {
	var rows []Domain
	var err error
	Dbo.Where("domain <> ?", "ALL").Find(&rows)
	return rows, err
}

// The ajax struct send as json
type DomainsPayload struct {
	Success bool     `json:"success"`
	Domains []Domain `json:"domains"`
	Error   string   `json:"error"`
}

// /domains
func HandleAjaxDomains(resp http.ResponseWriter, req *http.Request) {
	fmt.Println( "DOMAINSSS")
	payload := DomainsPayload{}
	payload.Success = true

	var err error
	payload.Domains, err = GetDomains()
	if err != nil {
		fmt.Println(err)
		payload.Error = "DB Error: " + err.Error()
	}

	base.WriteJSON(resp, payload)
}
