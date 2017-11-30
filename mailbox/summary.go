

package mailbox

import(

	"fmt"

	"bytes"

	"net/http"
	"net/mail"

	"github.com/mxk/go-imap/imap"


	"github.com/daffodil/go-postfixadmin/base"
)

//===================================================================
type Header struct {
	Uid uint32 `json:"uid"`

	FromName string `json:"from_name"`
	FromEmail string `json:"from_email"`
	Subject string `json:"subject"`

	Date string `json:"date"`

	Seen bool `json:"seen"`
	Flagged bool `json:"flagged"`
	Answered bool `json:"answered"`
}

type SummaryPayload struct {
	Headers[] *Header `json:"headers"`
	Success bool `json:"success"`
	Folders[] *Folder `json:"folders"`
	Uids[] uint32 `json:"uids"`

	Error string `json:"error"`
}


func HandleAjaxSummary(resp http.ResponseWriter, req *http.Request) {


	if base.AjaxAuth(resp, req) == false {
		return
	}

	client := CreateImapClient(resp, req)
	if client == nil {
		return
	}
	defer client.Logout(0)

	payload := new(SummaryPayload)
	payload.Folders = make([]*Folder, 0)
	payload.Uids = make([]uint32, 0)
	payload.Headers = make([]*Header, 0)
	payload.Success = true

	var err error
	payload.Folders, err = GetFolders(client)
	if err != nil {
		payload.Error = payload.Error + err.Error() + "\n"
	}

	//= Select inbox
	payload.Uids, err = GetUIDs("INBOX", client)


	uidlist := GetLastUIDs(payload.Uids)

	//----------------------------------------------
	//== Fetch last few  messages
	cmd, err := imap.Wait( client.UIDFetch(uidlist, "FLAGS", "INTERNALDATE", "RFC822.SIZE", "RFC822.HEADER") )
	if err != nil {
		fmt.Println("#################", err)
	}

	for _, rsp  := range cmd.Data {

		header := imap.AsBytes(rsp.MessageInfo().Attrs["RFC822.HEADER"])
		mm := new(Header)
		mm.Uid = rsp.MessageInfo().UID


		for flag, boo := range  rsp.MessageInfo().Flags {
			//fmt.Println( boo, flag, flag == "\\Seen" )
			if flag == "\\Seen" && boo {
				mm.Seen = true
			}

			if flag == "\\Flagged" && boo {
				mm.Flagged = true
			}
		}

		//fmt.Println("--------------------" )
		//fmt.Println("BS", imap.TypeOf( rsp.MessageInfo().Attrs["BODYSTRUCTURE"]) )
		//fmt.Println("BS",  rsp.MessageInfo().Attrs["BODYSTRUCTURE"] )

		for _, bsv := range imap.AsList(rsp.MessageInfo().Attrs["BODYSTRUCTURE"]) {
			fmt.Println("bsv=", imap.TypeOf( bsv ) )
			if  imap.TypeOf(bsv) == imap.List {
				for _, vx := range imap.AsList(bsv) {
					fmt.Println("  >", imap.TypeOf(vx), vx)
				}
				vvv := imap.AsList(bsv)
				fmt.Println(" ==", vvv[0], imap.TypeOf(vvv), vvv[0] == "application")
				if  imap.TypeOf(vvv[0]) == imap.QuotedString && imap.AsString(vvv[0]) == "application" {
					fmt.Println("@@@@@@@@@@@@", vvv[1])
				}

			}
		}

		if msg, _ := mail.ReadMessage(bytes.NewReader(header)); msg != nil {

			//fmt.Println(" msh_err= ", msg_err)
			mm.Subject = msg.Header.Get("Subject")


			// From
			from, fro_err := mail.ParseAddress(msg.Header.Get("From"))
			if fro_err != nil {
				fmt.Println("address ettot")
			} else {
				mm.FromName = from.Name
				mm.FromEmail = from.Address
			}

			// Date
			dat := imap.AsDateTime(rsp.MessageInfo().Attrs["INTERNALDATE"] )
			//if dat_err != nil {
			//	fmt.Println("date error", dat_err)
			//} else {
			mm.Date = dat.Format("2006-01-02 15:04:05")
			//}


			payload.Headers = append(payload.Headers, mm)
		}
		//mm := cmd.Data[midx].MessageInfo()


		//fmt.Println(mm.Attrs, mm.InternalDate)

		//payload.Messages = append(payload.Messages, mess)
		//}
	}


	base.SendPayload(resp, payload)


}

