

package sendmail

import(

	"fmt"
	"net/http"

	"time"

	"encoding/json"
)


type SendPayload struct {
	Success bool ` json:"success" `
	Message string ` json:"message" `
	Error string ` json:"error" `
}

func HandleAjaxSendTest(resp http.ResponseWriter, req *http.Request) {

	payload := SendPayload{}
	payload.Success = true

	// Set Ajax Headers
	resp.Header().Set("Content-Type", "application/json")

	tim := time.Now()
	t := tim.String()

	mess := Message{}
	mess.AddTo( Conf.TestEmail)
	mess.From = "noreply@" + Conf.DefaultDomain
	mess.Subject = "Test Message: " + t
	mess.Body = "Test" + t

	err := SendMessage(mess)
	if err != nil {
		payload.Error = err.Error()
	}

	json_str, _ := json.MarshalIndent(payload, "" , "  ")
	fmt.Fprint(resp, string(json_str))
}
