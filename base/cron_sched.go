
package base

import (

	"fmt"
	"time"
	"net/http"
	"reflect"

	"github.com/robfig/cron"
)

// Pointer to Cron instance..
// TODO: not yet implemented.. play code
var Cron *cron.Cron

// Starts the cron service
func StartCron() {

	Cron = cron.New()
	//Cron.AddFunc("@every 5s", CronTest )
	Cron.Start()
	fmt.Println("Cron Started" )
}

// Tests the cron.. wtf ???
func CronTest(){
	fmt.Println("CronTest " + NowStr())
	time.Sleep(3 * time.Second)
	fmt.Println("  >>>>>. " + NowStr())
}


// Struct to return ajax payload on current cronts and state from memeory..
type CronPayload struct {
	Success bool `json:"success"`
	Aliases string `json:"aliases"`
	Error string `json:"error"`
}



// handler for the /ajax?? /cron entries
// TODO
func HandleAjaxCron(resp http.ResponseWriter, req *http.Request) {

	if AjaxAuth(resp, req) == false {
		return
	}

	pay := CronPayload{}

	//fmt.Println( inspect(Cron.Entries()))
	for _, e := range Cron.Entries() {
		//st := reflect.TypeOf(e)
		fmt.Println(e.Next, reflect.TypeOf(e.Job))
	}

	SendPayload(resp, pay)


}