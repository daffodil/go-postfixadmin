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
func Initialize(conff *base.Config, db *sql.DB) {

	Conf = conff

	// This is bummer cos I want to use db.Driver.Name or alike instead of a new function var
	var err error
	Dbo, err = gorm.Open(Conf.Db.Engine, db)
	if err != nil {

	}
	Dbo.SingularTable(true)
	Dbo.LogMode(Conf.Debug)

	// TODO OOps a daisy.. in the nettles
	//base.Cron.AddFunc("@every 5s", VacationsExpire)

}
