package postfixadmin

import (
//"crypto"

)

// TODO:   Rip postfix auth in golang
// https://github.com/postfixadmin/postfixadmin/blob/master/functions.inc.php#L843
// vs http://golang.org/pkg/crypto/

func PassCrypt(raw_pass string, enc_pass string) (encrypted_pass string, err error) {

	if raw_pass != "" && enc_pass != "" {
		// had to use the vars ..
	}

	password := ""
	salt := ""

	if password != "" && salt != "" {
		// had to use the vars ..
	}

	//if crypto.MD5SHA1 == Conf.Crypto
	// WTF said pedro.. !!!

	return "", nil

}
