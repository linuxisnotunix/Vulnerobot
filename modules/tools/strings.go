package tools

import (
	"strings"
	"time"
)

//MonthsFr List french month
var MonthsFr = []string{"janvier", "février", "mars", "avril", "mai", "juin", "juillet", "août", "septembre", "octobre", "novembre", "décembre"}
var months = []time.Month{time.January, time.February, time.March, time.April, time.May, time.June, time.July, time.August, time.September, time.October, time.November, time.December}

//GetFrMonth parse a string to return the time.Month corresponding
func GetFrMonth(m string) time.Month {
	for i, mfr := range MonthsFr {
		if strings.Compare(m, mfr) == 0 {
			return months[i]
		}
	}
	return 0
}
