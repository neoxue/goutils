package date

import "time"

var DateDetailFormat = "2006-01-02 15:04:05"
var SimpleDate = "2006-01-02"
var DateForFile = "2006-01-02-15-04-05"

func Sprint(date time.Time, format string) string {
	return date.Format(format)
}
