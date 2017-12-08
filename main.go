package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
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
	data_source := config.Db.User + ":" + config.Db.Password + "@" + config.Db.Server + "/" + config.Db.Database
	var Db *sql.DB
	var err_db error
	Db, err_db = sql.Open(config.Db.Engine, data_source)
	if err_db != nil {
		fmt.Printf("Db Login Failed: ", err_db, "=", config.Db.Engine, data_source)
		os.Exit(1)
	}
	err_ping := Db.Ping()
	if err_ping != nil {
		fmt.Printf("Db Ping Failed: ", err_ping, "=", config.Db.Engine, data_source)
		os.Exit(1)
	}
	defer Db.Close()

	// Initialize modules
	base.Initialize(config)
	mailbox.Initialize(config)
	postfixadmin.Initialize(config, Db)
	sendmail.Initialize(config)


	// Main router
	BASE := "/api/v1"
	router := mux.NewRouter()

	//= Base
	router.HandleFunc(BASE, base.HandleAjaxInfo)


	//= Postfixadmin
	//pfaRouter := mux.NewRouter().PathPrefix(BASE + "/admin").Subrouter().StrictSlash(true)
	pfaRouter := router.PathPrefix("/admin").Subrouter()

	pfaRouter.HandleFunc("/domains", postfixadmin.HandleAjaxDomains)
	pfaRouter.HandleFunc("/domain/{domain}", postfixadmin.HandleAjaxDomain)
	pfaRouter.HandleFunc("/domain/{domain}/all", postfixadmin.HandleAjaxDomainAll)
	pfaRouter.HandleFunc("/domain/{domain}/vacations", postfixadmin.HandleAjaxVacations)
	pfaRouter.HandleFunc("/domain/{domain}/mailboxes", postfixadmin.HandleAjaxMailboxes)
	pfaRouter.HandleFunc("/domain/{domain}/virtual", postfixadmin.HandleAjaxDomainVirtual)

	pfaRouter.HandleFunc("/domain/{domain}/mailbox/{username}", postfixadmin.HandleAjaxMailbox)

	pfaRouter.HandleFunc("/mailbox/{email}", postfixadmin.HandleAjaxMailbox)
	pfaRouter.HandleFunc("/mailbox/{email}/vacation", postfixadmin.HandleAjaxVacation)
	pfaRouter.HandleFunc("/mailbox/{email}/set_secret", postfixadmin.HandleAjaxSetMailboxPassword)
	pfaRouter.HandleFunc("/mailbox/{email}/send_test", postfixadmin.HandleAjaxMailboxSendTest)

	pfaRouter.HandleFunc("/alias/{email}", postfixadmin.HandleAjaxAlias)
	pfaRouter.HandleFunc("/domain/{domain}/aliases", postfixadmin.HandleAjaxAliases)

	pfaRouter.HandleFunc("/vacation/notifications", postfixadmin.HandleAjaxVacationNotifications)

	pfaRouter.HandleFunc("/api/v1/admin/cron", base.HandleAjaxCron)


	//= SendMail
	router.HandleFunc("/api/v1/smtp/send_test", sendmail.HandleAjaxSendTest)

	//= Mailbox
	router.HandleFunc("/api/v1/mailbox/{address}/summary", mailbox.HandleAjaxSummary)
	router.HandleFunc("/api/v1/mailbox/{address}/folders", mailbox.HandleAjaxFolders)
	router.HandleFunc("/api/v1/mailbox/{address}/message/{folder}/{uid}", mailbox.AjaxMessageHandler)

	//router.PathPrefix(BASE + "/admin").Handler(negroni.New(
	//	Middleware1,
	//	Middleware2,
	//	negroni.Wrap(subRouter),
	//))
	common_midware := negroni.New()
	common_midware.UseHandler(router)

	//router.PathPrefix(BASE + "/admin").Handler(common_midware.With(
		//APIMiddleware1,
	//	negroni.Wrap(pfaRouter),
	//))
	// Start Http Server
	fmt.Println("Serving on " + config.HTTPListen)
	//http.Handle("/", router)
	//http.ListenAndServe(config.HTTPListen, n)
	common_midware.Run(config.HTTPListen)

}
