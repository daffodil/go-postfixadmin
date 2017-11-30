
package mailbox

import(

	"net/http"
	"github.com/mxk/go-imap/imap"

	"github.com/daffodil/go-postfixadmin/base"
)

//===================================================================

type Folder struct {
	Name string 		`json:"name"`
	Type string 		`json:"type"`
	HasChildren bool	`json:"has_children"`
	Unseen uint32			`json:"unseen"`
	Messages uint32		`json:"messages"`
	//Recent uint32			`json:"recent"`
}



//= Return list of IMAP folders
func GetFolders(client *imap.Client )([]*Folder, error){

	folders := make([]*Folder, 0)

	cmd, err := imap.Wait( client.List("", "*") )
	if err != nil {
		return folders, err
	}

	for idx := range cmd.Data {
		info := cmd.Data[idx].MailboxInfo()
		fol := new(Folder)
		fol.Name = info.Name

		for flag, boo := range  info.Attrs {
			//fmt.Println( info.Name, boo, flag)

			if info.Name == "INBOX"  && boo {
				fol.Type = "inbox"

			} else if flag == "\\Junk" && boo {
				fol.Type = "junk"

			} else if flag == "\\Trash" && boo {
				fol.Type = "trash"

			} else if flag == "\\Sent" && boo {
				fol.Type = "sent"

			} else if flag == "\\Drafts" && boo {
				fol.Type = "drafts"

			} else if flag == "\\Haschildren" && boo {
				fol.HasChildren = true

			} else if flag == "\\Hasnochildren" && boo {
				//

			} else {
				fol.Type = "??" + flag

			}

		}

		cmd_sta, err_sta := imap.Wait( client.Status(info.Name, "MESSAGES", "UNSEEN", "RECENT") )
		if err_sta != nil {
			//fmt.Println("STAT.ERR=", err_sta)
		} else {
			//fmt.Println("STAT.CMD=", cmd_sta.Data[0].MailboxStatus())
			ms := cmd_sta.Data[0].MailboxStatus()
			fol.Unseen = ms.Unseen
			fol.Messages = ms.Messages
			//fol.Recent = ms.Recent
		}

		folders = append(folders, fol)
	}

	return folders, nil
}


type FoldersPayload struct {
	Success bool `json:"success"`
	Folders[] *Folder `json:"folders"`
	Error string `json:"error"`
}

// /mailbox/<email>/folders
func HandleAjaxFolders(resp http.ResponseWriter, req *http.Request) {

	if base.AjaxAuth(resp, req) == false {
		return
	}

	client := CreateImapClient(resp, req)
	if client == nil {
		return
	}
	defer func() { client.Logout(0) }()

	payload := new(FoldersPayload)
	payload.Success = true

	var err error
	payload.Folders, err = GetFolders(client)
	if err != nil {
		payload.Error = err.Error()
	}

	base.SendPayload(resp, payload)
	//json_str, _ := json.MarshalIndent(payload, "" , "  ")
	//fmt.Fprint(resp, string(json_str))
}
