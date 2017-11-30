
package mailbox



import(


	"crypto/tls"

	"github.com/daffodil/go-postfixadmin/base"
)

//===============================================
type Cred  struct{
	Email string
	Password string
	Active uint
}

var conf *base.Config
var tlsConfig  *tls.Config


func Initialize( conff *base.Config){

	conf = conff

	tlsConfig = new(tls.Config)
	tlsConfig.ServerName = conf.ImapServer
	tlsConfig.InsecureSkipVerify = true

}
