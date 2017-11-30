package postfixadmin

import (
	"fmt"
	"strings"
	"time"

	"github.com/daffodil/go-postfixadmin/base"
	"github.com/daffodil/go-postfixadmin/sendmail"
)

var vacation_running bool

func VacationsExpire() {

	fmt.Println("CRON: VacationsExpire()")
	if vacation_running == true {
		fmt.Println("CRON: << already running VacationsExpire()")
		return
	}
	vacation_running = true
	started := time.Now()

	var active []string
	var upcoming []string
	var dead []string

	var reply string
	vacations, _ := GetVacations(Conf.DefaultDomain)
	//fmt.Println(vacations, err)

	//var to_delete []string

	for _, v := range vacations {
		now := time.Now()
		fmt.Println(v.Email, v.Activefrom, v.Activeuntil, now)

		// Handle Null Data
		if v.Activefrom == base.NULL_DATE_STR && v.Activeuntil == base.NULL_DATE_STR {

			reply = DeleteVacation("Null Data", v)
			if reply != "" {
				dead = append(dead, reply)
			}

		} else {

			// parse dates
			d_from, errf := base.ToDate(v.Activefrom)
			d_to, errt := base.ToDate(v.Activeuntil)

			//fmt.Println(errf, errt)
			if errf == nil && errt == nil {

				//fmt.Println(d_from, d_to)

				//expired := d_to.Before(now)
				if d_to.Before(now) {
					reply = DeleteVacation("Expired", v)
					dead = append(dead, reply)
					//to_delete = append(to_delete, v.Email)

				} else if d_from.After(now) {
					upcoming = append(upcoming, VacationToText(v))

				} else if d_to.After(now) {
					bod := VacationToText(v)
					if v.LastBot == "" {
						bod += SendVacationNotice(true, v)
					}
					active = append(active, bod)
				}
				//DeleteVacation(v)

			}
		}
	}
	//line := strings.Repeat("-", 30)

	bod := "" // "Started: " + base.NiceDT(started) + "\n"

	bod += makeSection("Active", active)
	bod += makeSection("Upcoming", upcoming)
	bod += makeSection("Dead", dead)

	//for _, em := range to_delete {

	//}

	/*
		bod += line
		bod += "# Active\n"
		bod += strings.Join(active, "\n")

		bod += line
		bod += "# Upcoming\n"
		bod += strings.Join(upcoming, "\n")

		bod += line
		bod += "# Dead\n"
		bod += strings.Join(dead, "\n")
	*/
	fmt.Println(base.NiceDateTime(started))
	fmt.Println(bod)

	sendi := false
	if sendi {
		sendmail.SendAdminMessage("Vacation Cron", bod)
	}
	vacation_running = false
}

func makeSection(h string, x []string) string {
	if len(x) == 0 {
		return ""
	}
	line := strings.Repeat("=", 30) + "\n"
	line += "# " + h + "\n"
	line += strings.Repeat("=", 30) + "\n\n"
	line += strings.Join(x, "\n")
	line += "\n"
	return line
}

func VacationToText(v *Vacation) string {
	m := "Email: " + v.Email + "\n"
	m += " From: " + v.Activefrom + "\n"
	m += "Until: " + v.Activeuntil + "\n"
	return m
}

func DeleteVacation(status string, v *Vacation) string {

	m := VacationToText(v)

	if true || Conf.Live == true {
		//Dbo.Delete(&v)
		m += "## Deleted ##\n"
		nos, err := GetVacationNotifications(v.Email, "email")
		if err == nil {
		}
		for _, r := range nos {
			dt, _ := base.ToDate(r.NotifiedAt)
			m += base.NiceDateTime(dt) + " " + r.Notified + "\n"
		}
		//bo.
	}

	return m
}

func SendVacationNotice(active bool, v *Vacation) string {

	s := v.Email + ": "
	if active {
		s += "Out Of Office: Activated "
	} else {
		s += "Out Of Office: Deactivated "
	}
	m := VacationToText(v)
	sendmail.SendAdminMessage(s, m)
	if active {
		v.LastBot = base.NowStr()
		Dbo.Save(v)
	}
	return s + "\n"

}
