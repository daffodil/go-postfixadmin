package postfixadmin

import (
	"database/sql"

	//"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"github.com/daffodil/go-postfixadmin/base"
)

var Conf *base.Config
var Dbo *gorm.DB

// Initializes the postfix admin module..
func Initialize(conf *base.Config, db *sql.DB) {

	Conf = conf

	// Create gorm instance
	var err error
	Dbo, err = gorm.Open(Conf.Db.Engine, db)
	if err != nil {

	}
	Dbo.SingularTable(true)
	Dbo.LogMode(Conf.Debug)

	// TODO Start cron
	//base.Cron.AddFunc("@every 5s", VacationsExpire)

}
