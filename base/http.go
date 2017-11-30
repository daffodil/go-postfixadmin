
package base

import(

	"fmt"
	"net/http"
	"encoding/json"
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
func AjaxAuth(resp http.ResponseWriter, req *http.Request) bool {

	// Set Ajax Headers ie were in json land
	resp.Header().Set("Content-Type", "application/json")


	switch req.Method {

	// A HTTP POST so check the auth secret
	case "POST":
		req.ParseForm()
		if req.Form.Get("auth") != conf.AuthSecret {
			http.Error(resp, CreatePermissionErrPayload(), 500)
			return false
		}

	// This is well dodgy here bill.. make it GET only as its assumed ???
	default:

		if req.URL.Query().Get("auth") != conf.AuthSecret {
			http.Error(resp, CreatePermissionErrPayload(), 500)
			return false
		}
	}

	return true
}

// Writes out the "dict/map" in json to remote http client
func SendPayload(resp http.ResponseWriter, payload interface{} ) {
	json_str, _ := json.MarshalIndent(payload, "" , "")
	fmt.Fprint(resp, string(json_str))
}

// Struct for sending an ajax error
type ErrorPayload struct {
	Success bool `json:"success"`
	Error string `json:"error"`
}

// Send the error payload json enoded to client..
func SendErrorPayload(resp http.ResponseWriter, err string){

	payload := new(ErrorPayload)
	payload.Success = true
	payload.Error = err

	SendPayload(resp, payload)
}