package timeutil

import "time"

func ParseLocaltime(t string) (local time.Time) {
	loc, _ := time.LoadLocation("Asia/Makassar")
	local, _ = time.ParseInLocation("2006-01-02 15:04:05", t, loc)
	return
}

func ConvertLocalTime(t time.Time) (local time.Time) {
	loc, _ := time.LoadLocation("Asia/Makassar")
	return t.In(loc)
}
