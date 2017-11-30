package postfixadmin

import (
	"errors"
	"net/http"

	"github.com/cenkalti/log"
	"github.com/gorilla/mux"

	"github.com/daffodil/go-postfixadmin/base"
)

type Domain struct {
	Domain      string ` json:"domain" gorm:"primary_key"`
	Description string ` json:"description"`
	Aliases     int    ` json:"aliases"`
	Mailboxes   int    ` json:"mailboxes"`
	MaxQuota    int    ` json:"maxquota"`
	Quota       int    ` json:"quota"`
	Transport   string ` json:"transport"`
	BackupMx    int    ` json:"backupmx"`
	Created     string ` json:"created"`
	Modified    string ` json:"modified"`
	Active      int    ` json:"active"`
}

// Table name for gorm
func (me Domain) TableName() string {
	return Conf.Db.TableNames["domain"]
}

// Checks if domain exist and active
func IsDomainValid(domain_name string) error {

	if DomainExists(domain_name) == false {
		return errors.New("Domain `" + domain_name + "' does not exist")
	}
	return nil
}

// Check domain record exists
func DomainExists(domain string) bool {
	var count int
	Dbo.Model(Domain{}).Where("domain = ?", domain).Count(&count)
	if count == 0 {
		return false
	}
	return true
}

//= Get domain row from db
func GetDomain(domain_name string) (Domain, error) {
	var dom Domain
	var err error

	err = IsDomainValid(domain_name)
	if err != nil {
		return dom, err
	}
	Dbo.Where("domain = ? ", domain_name).Order("domain").Find(&dom)
	return dom, err
}

//= Ajax struct for `domain`
type DomainPayload struct {
	Success bool   `json:"success"` // keep extjs happy
	Domain  Domain `json:"domain"`
	Error   string `json:"error"`
}

// /domain/{domain}
func HandleAjaxDomain(resp http.ResponseWriter, req *http.Request) {

	log.Info("AjaxHandleDomain")
	if base.AjaxAuth(resp, req) == false {
		return
	}
	vars := mux.Vars(req)
	// TODO check var is valid

	payload := DomainPayload{}
	payload.Success = true

	var err error
	payload.Domain, err = GetDomain(vars["domain"])
	if err != nil {
		log.Info(err.Error())
		payload.Error = "" + err.Error()
	}

	base.SendPayload(resp, payload)

}

// Ajax struct for `domain` all
type DomainAllPayload struct {
	Success   bool      `json:"success"` // keep extjs happy
	Domain    Domain    `json:"domain"`
	Mailboxes []Mailbox `json:"mailboxes"`
	Aliases   []Alias   `json:"aliases"`
	Error     string    `json:"error"`
}

//  /ajax/domain/{domain}/all
func HandleAjaxDomainAll(resp http.ResponseWriter, req *http.Request) {

	if base.AjaxAuth(resp, req) == false {
		return
	}
	log.Info("DomainAllAjaxHandler")

	vars := mux.Vars(req)
	domain := vars["domain"]

	payload := DomainAllPayload{}
	payload.Success = true

	var err error
	payload.Domain, err = GetDomain(domain)
	if err != nil {
		log.Info(err.Error())
		payload.Error = "" + err.Error()
	}
	payload.Mailboxes, err = GetMailboxes(domain)
	if err != nil {
		log.Info(err.Error())
		payload.Error = "" + err.Error()
	}
	payload.Aliases, err = GetAliases(domain)
	if err != nil {
		log.Info(err.Error())
		payload.Error = "" + err.Error()
	}

	base.SendPayload(resp, payload)
}
