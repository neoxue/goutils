package date

import "time"

var SimpleDateFormat = "2006-01-02 15:04:05"

func Sprint(date time.Time, format string) string {
	return date.Format(format)
}
