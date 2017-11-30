package postfixadmin

import (
	"fmt"
	"net/http"
	"errors"

	"github.com/gorilla/mux"

	"github.com/daffodil/go-postfixadmin/base"
)

type AliasesPayload struct {
	Success bool    `json:"success"` // keep extjs happy
	Aliases []Alias `json:"aliases"`
	Error   string  `json:"error"`
}

func GetAliases(domain string) ([]Alias, error) {

	var rows []Alias
	if DomainExists(domain) == false {
		return rows, errors.New("Domain '" + domain + "` does not exist")
	}
	var err error
	Dbo.Where("domain=?", domain).Order("address").Find(&rows)
	return rows, err
}

// /domain/<domain>/aliases
func HandleAjaxAliases(resp http.ResponseWriter, req *http.Request) {

	if base.AjaxAuth(resp, req) == false {
		return
	}
	fmt.Println("AliasesAjaxHandler")
	vars := mux.Vars(req)

	payload := AliasesPayload{}
	payload.Success = true
	payload.Aliases = make([]Alias, 0)

	var err error
	payload.Aliases, err = GetAliases(vars["domain"])
	if err != nil {
		fmt.Println(err)
		payload.Error = "Error: " + err.Error()
	}

	base.SendPayload(resp, payload)
}
