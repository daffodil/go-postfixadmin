

package base

import(

)

// Main config pointer
var conf *Config



// The main `config.yaml` reader/struct
type Config struct {

	Debug bool `yaml:"debug" json:"debug" `
	Live bool `yaml:"live" json:"live" `

	HTTPListen string `yaml:"http_listen" json:"http_listen"`

	// If true with return HTTP status error codes, in addition to json.
	// Eg an authentication error of 304, in addition to the json data error
	// Some application would want to have 200 ok, and the error in json
	HTTPErrors bool `yaml:"http_errors" json:"http_errors"`
	JSONPretty bool `yaml:"json_pretty" json:"json_pretty"`

	Tokens []TokenConf `yaml:"tokens" json:"tokens" `

	EmailPrefix string `yaml:"email_prefix" json:"email_prefix" `
	AdminEmail string `yaml:"admin_email" json:"admin_email" `
	FromEmail string `yaml:"from_email" json:"from_email" `
	SyslogEmail string `yaml:"syslog_email" json:"syslog_email" `
	TestEmail string `yaml:"test_email" json:"test_email" `
	NoReplyEmail string `yaml:"noreply_email" json:"noreply_email" `

	Db DbConf

	DefaultDomain string `yaml:"default_domain" json:"default_domain" `
	VacationDomain string `yaml:"vacation_domain" json:"vacation_domain" `




	SMTPServer SMTPConf `yaml:"smtp" json:"smtp"`

	IMAPServer IMAPConf `yaml:"imap" json:"imap"`
}


// Database config (TODO: xtra params, ssl et all)
type DbConf struct {
	Engine string 	`yaml:"engine" json:"engine" `
	Server string 	`yaml:"server" json:"server"`
	User string 	`yaml:"user" json:"user"`
	Password string `yaml:"password" json:"password"`
	Database string `yaml:"database" json:"database"`
	Debug bool 		`yaml:"debug" json:"debug"`
	TableNames map[string]string  `yaml:"table_names" json:"table_names"`
}

// SMTP server config
type SMTPConf struct {
	Server string 	` yaml:"server" json:"server" `
	Port int 		` yaml:"port" json:"port" `
	Login string 	`yaml:"login" json:"login"`
	Password string 	`yaml:"password" json:"password"`
	TestMode bool `yaml:"test_mode" json:"test_mode"`
	//TestEmail string `yaml:"test_email" json:"test_email"`
}

// IMAP server
type IMAPConf struct {
	Active bool 	` yaml:"active" json:"active" `
	Server string  	` yaml:"server" json:"server" `
}


// Simple TOKEN auth
type TokenConf struct {
	Active bool 	` yaml:"active" json:"active" `
	Secret string  	` yaml:"secret" json:"secret" `
	Header string 	` yaml:"header" json:"header"`
	Ips []string 	` yaml:"ips" json:"ips"`
}
