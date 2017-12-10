

package mailbox

import(

	"fmt"

	"bytes"

	"net/http"
	"net/mail"
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/mxk/go-imap/imap"
	"github.com/jhillyerd/go.enmime"

)

type XMeta struct {
	batch_id int ` json:"batch_id" `
	job_no string ` json:"job_no" `
}

func ParseMeta(s string ) *XMeta {

	if s == "" {
		return nil
	}
	m := new(XMeta)
	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		return nil
	}
	return m
}

type Attachment struct {
	FileName string  ` json:"file_name" `
	Type string  ` json:"type" `
	Blob []byte  ` json:"blob" `
}

//===================================================================
type Message struct {
	Uid uint32 `json:"uid"`
	Folder string `json:"folder"`

	FromName string `json:"from_name"`
	FromEmail string `json:"from_email"`

	Date string `json:"date"`

	Seen bool `json:"seen"`
	Flagged bool `json:"flagged"`
	Answered bool `json:"answered"`

	Subject string `json:"subject"`
	BodyHtml string  `json:"body_html"`
	BodyText string  `json:"body_text"`

	ContentType string `json:"content_type"`

	XMeta *XMeta `json:"x_meta"`
	Attachments []Attachment `json:"attachments"`
}




func GetMessage(folder, uid string, client *imap.Client) (messag *Message, e error ){

	cmd, err := client.Select(folder, true)
	if err != nil {
		return nil, err
	}

	uidlist, _ := imap.NewSeqSet(uid)
	//uidlist.Add(uid)

	fmt.Println("get_mess", folder, uid)
	mess := new(Message)
	mess.Folder = folder

	cmd, err = imap.Wait( client.UIDFetch(uidlist, "FLAGS", "INTERNALDATE", "RFC822.SIZE",  "RFC822")) //,  "RFC822.HEADER")) //,  "BODY.PEEK[TEXT]") )
	if err != nil {
		return mess, err
	}
	fmt.Println("Data=", len(cmd.Data), cmd.Data)
	rsp := cmd.Data[0]
	minfo := rsp.MessageInfo()
	fmt.Println("Info=", minfo.UID, minfo.Seq, minfo.Size, minfo.InternalDate)
	mess.Uid = minfo.UID


	fmt.Println("Flags")
	for flag, boo := range  minfo.Flags {
		fmt.Println("=", flag, boo)
		if flag == "\\Seen" && boo {
			mess.Seen = true
		}
		if flag == "\\Flagged" && boo {
			mess.Flagged = true
		}
	}
	fmt.Println("Attrs")
	for ki := range  minfo.Attrs {
		fmt.Println("=", ki)
	}

	// Date
	dat := imap.AsDateTime(minfo.Attrs["INTERNALDATE"])
	mess.Date = dat.Format("2006-01-02 15:04:05")

	//fmt.Println("dt==", minfo.Attrs["INTERNALDATE"])
	//fmt.Println("head==", imap.TypeOf( minfo.Attrs["RFC822.HEADER"]) )
	//fmt.Println("head==", minfo.Attrs["RFC822.HEADER"] )
	/*
	bites := imap.AsBytes(minfo.Attrs["RFC822"])
	msg, msg_err := mail.ReadMessage(bytes.NewReader(bites))
	if msg_err != nil {
		return mess, msg_err
	}
	*/
	//fmt.Println("@@@@@", string(imap.AsBytes(minfo.Attrs["BODYSTRUCTURE"])))
	msg, errrm := mail.ReadMessage( bytes.NewReader(imap.AsBytes(minfo.Attrs["RFC822"]) ))
	if errrm != nil {
		fmt.Println("errrm=", errrm)
	}
	mime, mime_err := enmime.ParseMIMEBody(msg)
	if mime_err != nil {
		fmt.Println("err=", mime_err, mime)
	}

	// From
	from, fro_err := mail.ParseAddress(msg.Header.Get("From"))
	if fro_err != nil {
		fmt.Println("address ettot")
	} else {
		mess.FromName = from.Name
		mess.FromEmail = from.Address
	}

	mess.Subject = msg.Header.Get("Subject")
	mess.ContentType = msg.Header.Get("Content-Type")

	dd, _ := msg.Header.Date()
	fmt.Println("DAT=", mess.Date,dd )
	//fmt.Println("body=", cmd.Data[0].String)

	mess.BodyText =  mime.Text
	mess.BodyHtml =  mime.HTML //imap.AsString(minfo.Attrs["RFC822"])

	// Meta Data
	xmeta := msg.Header.Get("X-gStl-META")
	if xmeta != "" {
		mess.XMeta = ParseMeta(xmeta)
	}
	//fmt.Println("META=", xmeta, mess.XMeta)

	lst, _ := msg.Header.AddressList("to")
	fmt.Println("to=", lst)

	lst2, _ := msg.Header.AddressList("cc")
	fmt.Println("META=", lst2)


	// Attacjments
	fmt.Println("attach=", len(mime.Attachments))
	for i, a := range mime.Attachments {
		att := Attachment{FileName: a.FileName(), Type: a.ContentType(), Blob: a.Content()}
		mess.Attachments = append(mess.Attachments, att)
		fmt.Println("==", i, a.FileName())
	}

	return mess, nil

}

type MessagePayload struct {

	Success bool `json:"success"`
	Message *Message `json:"message"`

	Error string `json:"error"`
}

func AjaxMessageHandler(resp http.ResponseWriter, req *http.Request) {


	client := CreateImapClient(resp, req)
	if client == nil {
		return
	}
	defer func() { client.Logout(0) }()

	payload := new(MessagePayload)
	payload.Success = true
	payload.Message = new(Message)

	var err error
	vars := mux.Vars(req)
	payload.Message, err = GetMessage(vars["folder"], vars["uid"], client)
	if err != nil {
		fmt.Println("err", err)

	}
	//payload.Message.Folder = vars["folder"]
	//payload.Message.Uid = vars["uid"]

	json_str, _ := json.MarshalIndent(payload, "" , "  ")
	fmt.Fprint(resp, string(json_str))

}



