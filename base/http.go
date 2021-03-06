
package base

import(

	"fmt"
	"net/http"
	"encoding/json"

)

const(
	V1_BASE = "/api/v1"
)

// Return ajax struct created in error,,
type ErrPayload struct {
	Success bool `json:"success"` // keep extjs happy
	Error string `json:"error"`
}

// Serialises a Permission payload to json
func CreatePermissionErrPayload() string {

	payload := ErrPayload{}
	payload.Success = true
	payload.Error = "Permission denied"
	json_str, _ := json.MarshalIndent(payload, "" , "  ")
	return string(json_str)
}

// This is atmo a simple auth mechanism
// and a pass down the line function shite.. which golang does..
func DEADAjaxAuth(resp http.ResponseWriter, req *http.Request) bool {

	// Set Ajax Headers ie were in json land
	resp.Header().Set("Content-Type", "application/json")

	// simple token auth
	/*
	if conf.Token.Active {
		token := req.Header.Get(conf.Token.Header)
		if len(token) > 5 && token != conf.Token.Secret {

			// check ip match
			real_ip := req.Header.Get("X-Real-IP")
			for _, v := range conf.Token.Ips {
				if v == real_ip {
					return true
				}
			}
		}
	} */
	resp.WriteHeader(http.StatusUnauthorized)
	resp.Write([]byte("500 - postfixadmin permission error"))
	return false
}

// Writes out the payload in json
func WriteJSON(resp http.ResponseWriter, payload interface{} ) {
	indent := ""
	if conf.JSONPretty {
		indent = "  "
	}
	json_str, _ := json.MarshalIndent(payload, "" , indent)
	fmt.Fprint(resp, string(json_str))
}

// Struct for sending an ajax error
type ErrorPayload struct {
	Success bool `json:"success"`
	Error string `json:"error"`
}

// Send the error payload json enoded to client..
func WriteErrorJSON(resp http.ResponseWriter, err string){

	payload := new(ErrorPayload)
	payload.Success = true
	payload.Error = err

	WriteJSON(resp, payload)
}


//  /api/v1 info
func HandleAjaxInfo(resp http.ResponseWriter, req *http.Request) {

	pay := make(map[string]string)

	pay["real_ip"] = req.Header.Get("X-Real-IP")
	pay["remote_addr"] = req.RemoteAddr

	WriteJSON(resp, pay)


}
