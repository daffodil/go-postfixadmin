package mailbox


import(

	"github.com/mxk/go-imap/imap"

)



func GetUIDs(mbox string, client *imap.Client) ([]uint32, error ) {

	uids := make([]uint32, 0)

	cmd, err := client.Select(mbox, true)
	if err != nil {
		return uids, err
	}

	//== Get UIDS of all messages
	cmd, err = imap.Wait( client.UIDSearch("", "ALL") )
	if err != nil {
		return uids, err
	}

	for idx := range cmd.Data {
		for _, uid := range cmd.Data[idx].SearchResults() {
			uids = append(uids, uid)
		}
	}
	return uids, nil

}
func GetLastUIDs(alluids []uint32) *imap.SeqSet {
	//payload.Uids, err = GetUIDs("INBOX", client)

	//= Calc last few messages
	lenny := len(alluids)
	last := lenny - 50  // ################
	if  last < 0 {
		last = 0
	}
	//= Make List of messages uids
	uidlist, _ := imap.NewSeqSet("")
	for _, uid := range alluids[last:lenny] {
		uidlist.AddNum(uid)
	}
	return uidlist
}
