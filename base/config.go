

package base

import(

)

// Main config pointer
var conf *Config

type DbConf struct {
	Engine string 	`yaml:"engine" json:"engine" `
	Server string 	`yaml:"server" json:"server"`
	User string 	`yaml:"user" json:"user"`
	Password string `yaml:"password" json:"password"`
	Database string `yaml:"database" json:"database"`
	Debug bool 		`yaml:"debug" json:"debug"`
	TableNames map[string]string  `yaml:"table_names" json:"table_names"`
}

type SMTPConf struct {
	Server string 	` yaml:"server" json:"server" `
	Port int 		` yaml:"port" json:"port" `
	Login string 	`yaml:"login" json:"login"`
	Password string 	`yaml:"password" json:"password"`
	TestMode bool `yaml:"test_mode" json:"test_mode"`
	//TestEmail string `yaml:"test_email" json:"test_email"`
}

// The main `config.yaml` reader struct where yaml is serialed into
type Config struct {

	Debug bool `yaml:"debug" json:"debug" `
	Live bool `yaml:"live" json:"live" `

	AuthSecret string `yaml:"auth_secret" json:"auth_secret" `

	EmailPrefix string `yaml:"email_prefix" json:"email_prefix" `
	AdminEmail string `yaml:"admin_email" json:"admin_email" `
	FromEmail string `yaml:"from_email" json:"from_email" `
	SyslogEmail string `yaml:"syslog_email" json:"syslog_email" `
	TestEmail string `yaml:"test_email" json:"test_email" `
	NoReplyEmail string `yaml:"noreply_email" json:"noreply_email" `

	Db DbConf

	DefaultDomain string `yaml:"default_domain" json:"default_domain" `
	VacationDomain string `yaml:"vacation_domain" json:"vacation_domain" `

	HTTPListen string `yaml:"http_listen" json:"http_listen"`
	//IMAPAddress string `yaml:"imap_adddress" json:"imap_adddress"`

	SMTPLogin SMTPConf `yaml:"smtp" json:"smtp"`

	ImapServer string `yaml:"imap_server" json:"imap_server"`
}






