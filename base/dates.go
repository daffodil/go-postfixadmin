
package base

import(

	"time"

)

const (
// The database struct may return this, so its a NULL date as a string compare
	NULL_DATE_STR = "0000-00-00 00:00:00"
	MYSQL_DATE_TIME = "2006-01-02 15:04:05"

)



// returns current time as a string for mysql..
// TODO consider tz and UTC for future..
func NowStr() string {

	t := time.Now()
	return t.Format("2006-01-02 15:04:05")
}





// Converts a mysql string data to a time instance
func ToDate(s string) (time.Time, error) {
	return time.Parse(MYSQL_DATE_TIME, s)

}

// Prints date in a nice readable  format.
func NiceDateTime(t time.Time) string {
	return t.Format(time.RFC850)
}
