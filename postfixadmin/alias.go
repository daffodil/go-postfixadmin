package postfixadmin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/daffodil/go-postfixadmin/base"
)

// Represents the alias table
type Alias struct {
	Address  string `json:"address" gorm:"primary_key" `
	Goto     string `json:"goto"`
	Domain   string `json:"domain"`
	Created  string `json:"created"`
	Modified string `json:"modified"`
	Active   int    `json:"active"`
}

func (me *Alias) TableName() string {
	return Conf.Db.TableNames["alias"]
}

// Save instance to Db
func (me *Alias) Save() {
	me.Modified = base.NowStr()
	Dbo.Save(&me)
}

// TODO make this not remove vacation ?
func (me *Alias) ClearAllGoto() {
	me.Goto = ""
}

// Add address to forwarding
func (me *Alias) AddGoto(addr string) {

	var parts []string
	if me.Goto == "" {
		parts = make([]string, 0)
		parts = append(parts, addr)

	} else {
		parts = strings.Split(me.Goto, ",")
		//fmt.Println("Pre-parts=", len(parts), parts,)
		found := false
		for _, p := range parts {
			if p == addr {
				found = true
			}
		}
		if found == true {
			//fmt.Println("DOun vac alias")
			return
		}
		parts = append(parts, addr)
	}
	me.Goto = strings.Join(parts, ",")
	//fmt.Println("Post-parts=", len(parts),  parts, me.Goto)
}

// Remove address from forwarding
func (me *Alias) RemoveGoto(addr string) {

	addresses := make([]string, 0)
	gotos := strings.Split(me.Goto, ",")

	for _, p := range gotos {
		if p != addr {
			addresses = append(addresses, p)
		}
	}
	me.Goto = strings.Join(addresses, ",")
}

func AliasExists(address string) bool {
	var count int
	Dbo.Model(Alias{}).Where("address = ?", address).Count(&count)
	if count == 0 {
		return false
	}
	return true
}

func GetAlias(email string) (*Alias, error) {

	alias := new(Alias)
	var err error
	Dbo.Where("address = ? ", email).Find(&alias)
	return alias, err
}

// The ajax payload  container
type AliasPayload struct {
	Success bool   `json:"success"` // keep extjs happy
	Alias   *Alias `json:"alias"`
	Error   string `json:"error"`
}

// Crreates a payload and struct.. defined above
func CreateAliasPayload() AliasPayload {
	payload := AliasPayload{}
	payload.Success = true
	//payload.Alias = make(Alias, 0)
	return payload
}

// /alias/<email>
func HandleAjaxAlias(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("AliasAjaxHandler")

	if base.AjaxAuth(resp, req) == false {
		return
	}

	payload := CreateAliasPayload()
	vars := mux.Vars(req)

	if vars["email"] == "new" {
		payload.Alias = new(Alias)
		base.SendPayload(resp, payload)
		return
	}

	addr, erra := ParseAddress(vars["email"])
	if erra != nil {
		payload.Error = erra.Error()
		base.SendPayload(resp, erra)
		return
	}

	var err error

	alias_exists := AliasExists(addr.Address)

	if alias_exists {
		payload.Alias, err = GetAlias(addr.Address)
		if err != nil {
			fmt.Println(err)
			payload.Error = "DB Error: " + err.Error()
		}
	} else {
		payload.Alias = new(Alias)
	}

	switch req.Method {

	case "POST":

		f := req.Form
		if alias_exists == false {
			payload.Alias.Address = addr.Address
			payload.Alias.Domain = addr.Domain
			Dbo.Create(payload.Alias)
		}
		payload.Alias.Goto = f.Get("goto")
		payload.Alias.Active = 1
		Dbo.Save(payload.Alias)

	}

	base.SendPayload(resp, payload)
}
