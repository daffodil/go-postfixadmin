

package sendmail

import (
	//"github.com/gorilla/mux"
	"github.com/daffodil/go-postfixadmin/base"
)
/*
smtp:
	server: localhost
	login: mash
	password: root
(
*/



var Conf *base.Config

func Initialize(conf_i *base.Config) {

	Conf = conf_i
}

