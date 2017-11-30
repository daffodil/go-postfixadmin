

package mailbox

import (
	"net/mail"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mxk/go-imap/imap"

	//"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/daffodil/go-postfixadmin/base"
	"github.com/daffodil/go-postfixadmin/postfixadmin"
)

//== Check and get credentials of email
func CreateImapClient(resp http.ResponseWriter, req *http.Request)(*imap.Client) {

	//= Get Email address and validate
	vars := mux.Vars(req)

	addr, addr_err := mail.ParseAddress(vars["address"])
	if addr_err != nil {
		base.SendErrorPayload(resp, "Invalid email address")
		return nil
	}

	// Get email and Password and active from DB and validate
	pass, err := postfixadmin.GetMailboxPassword(addr.Address)
	if err != nil {
		base.SendErrorPayload(resp, err.Error())
		return nil
	}

	//= Connect to Server
	client, conn_err := imap.DialTLS(conf.ImapServer, tlsConfig)
	if conn_err != nil {
		base.SendErrorPayload(resp, conn_err.Error())
		return nil
	}
	_, login_err := client.Login( addr.Address, pass )
	if login_err != nil {
		base.SendErrorPayload(resp, "IMAP login error")
		return nil
	}
	return client

}