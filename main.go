package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/daffodil/go-postfixadmin/base"
	"github.com/daffodil/go-postfixadmin/mailbox"
	"github.com/daffodil/go-postfixadmin/postfixadmin"
	"github.com/daffodil/go-postfixadmin/sendmail"
)

func main() {

	config_file := flag.String("config", "config.yaml", "Config file")
	flag.Parse()

	// Create and load config.yaml
	config := new(base.Config)
	contents, e := ioutil.ReadFile(*config_file)
	if e != nil {
		fmt.Printf("Config File Error: %v\n", e)
		fmt.Printf("create one with -w \n")
		os.Exit(1)
	}
	if err := yaml.Unmarshal(contents, &config); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	// Create Database connection + ping
	var Db *sql.DB
	var err_db error
	Db, err_db = sql.Open(config.Db.Engine, config.Db.Datasource)
	if err_db != nil {
		fmt.Printf("Db Login Failed: ", err_db, "=", config.Db.Engine, config.Db.Datasource)
		os.Exit(1)
	}
	err_ping := Db.Ping()
	if err_ping != nil {
		fmt.Printf("Db Ping Failed: ", err_ping, "=", config.Db.Engine, config.Db.Datasource)
		os.Exit(1)
	}
	defer Db.Close()

	// Initialize modules
	base.Initialize(config)
	mailbox.Initialize(config)
	postfixadmin.Initialize(config, Db)
	sendmail.Initialize(config)

	// Setup www router and config modules
	router := mux.NewRouter()

	//= Base
	router.HandleFunc("/api/v1/admin/cron", base.HandleAjaxCron)

	//= Postfixadmin
	router.HandleFunc("/api/v1/admin/domains", postfixadmin.HandleAjaxDomains)
	router.HandleFunc("/api/v1/admin/domain/{domain}", postfixadmin.HandleAjaxDomain)
	router.HandleFunc("/api/v1/admin/domain/{domain}/all", postfixadmin.HandleAjaxDomainAll)
	router.HandleFunc("/api/v1/admin/domain/{domain}/vacations", postfixadmin.HandleAjaxVacations)
	router.HandleFunc("/api/v1/admin/domain/{domain}/mailboxes", postfixadmin.HandleAjaxMailboxes)
	router.HandleFunc("/api/v1/admin/domain/{domain}/virtual", postfixadmin.HandleAjaxDomainVirtual)

	router.HandleFunc("/api/v1/admin/domain/{domain}/mailbox/{username}", postfixadmin.HandleAjaxMailbox)

	router.HandleFunc("/api/v1/admin/mailbox/{email}", postfixadmin.HandleAjaxMailbox)
	router.HandleFunc("/api/v1/admin/mailbox/{email}/vacation", postfixadmin.HandleAjaxVacation)
	router.HandleFunc("/api/v1/admin/mailbox/{email}/set_secret", postfixadmin.HandleAjaxSetMailboxPassword)
	router.HandleFunc("/api/v1/admin/mailbox/{email}/send_test", postfixadmin.HandleAjaxMailboxSendTest)

	router.HandleFunc("/api/v1/admin/alias/{email}", postfixadmin.HandleAjaxAlias)
	router.HandleFunc("/api/v1/admin/domain/{domain}/aliases", postfixadmin.HandleAjaxAliases)

	router.HandleFunc("/api/v1/admin/vacation/notifications", postfixadmin.HandleAjaxVacationNotifications)

	//= SendMail
	router.HandleFunc("/api/v1/smtp/send_test", sendmail.HandleAjaxSendTest)

	//= Mailbox
	router.HandleFunc("/api/v1/mailbox/{address}/summary", mailbox.HandleAjaxSummary)
	router.HandleFunc("/api/v1/mailbox/{address}/folders", mailbox.HandleAjaxFolders)
	router.HandleFunc("/api/v1/mailbox/{address}/message/{folder}/{uid}", mailbox.AjaxMessageHandler)

	// Start Http Server
	fmt.Println("Serving on " + config.HTTPListen)
	http.Handle("/", router)
	http.ListenAndServe(config.HTTPListen, nil)

}
