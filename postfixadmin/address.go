package postfixadmin

import (
	"errors"
	"strings"
)

//= Components of an email address (pedro hack)
type Address struct {
	Address         string // foo@example.com
	User            string // foo
	Domain          string // example.com
	VacationAddress string // foo#example.com@autoreply.example.com
}

// Parses an email_address to Addr{} or error
// Pedro hack to parse address in one place to constitient parts..
func ParseAddress(email_address string) (*Address, error) {

	stripped := strings.TrimSpace(email_address)
	if len(stripped) == 0 {
		return nil, errors.New("Invalid Email - zero length")
	}

	if strings.Contains(stripped, "@") == false {
		return nil, errors.New("Invalid Email, no @ in `" + email_address + "` ")
	}

	user_domain := strings.Split(stripped, "@")
	if DomainExists(user_domain[1]) == false {
		return nil, errors.New("Domain not exist in Db for email `" + email_address + "` ")
	}

	addr := new(Address)
	addr.Address = stripped
	addr.User = user_domain[0]
	addr.Domain = user_domain[1]
	addr.VacationAddress = addr.User + "#" + addr.Domain + "@" + Conf.VacationDomain

	return addr, nil

}
