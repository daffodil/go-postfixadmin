
package base

import(
	"net/http"
	"fmt"
)


func AuthMiddleware(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	// do some stuff before
	fmt.Println("at", req.URL.String())

	// Check tokens
	if len(conf.Tokens) > 0 {
		for _, t := range conf.Tokens {
			fmt.Println("t=", t)
			if t.Active {
				// check token matches
				token := req.Header.Get(t.Header)
				if len(token) > 5 && token ==  t.Secret {
					// check ip match
					real_ip := req.Header.Get("X-Real-IP")
					for _, v := range t.Ips {
						if v == real_ip {
							fmt.Println("YES", t.Header)
						}
					}
				}
			}
		}
	}
	fmt.Println("AUTH FAIL")
	if conf.HTTPErrors {
		http.Error(rw, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	}
	WriteErrorJSON(rw, "Authentication Failed")
	return
	fmt.Println("AUTH FAIL after")
	next(rw, req)
	// do some stuff after
}