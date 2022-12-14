package main

import (
	"time"
)

func parseDateStamp(s string) time.Time {
	pattern := "2006-01-02"
	t, _ := time.Parse(pattern, s)
	return t
}

// func incrementDay(t *time.Time) {
// *t = time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, time.UTC)
// 	*t = t.Add(time.Hour * 24)
// }

// dayIndex returns the index of the day in the first/last
// range. If out of range, returns -1
func dayIndex(given time.Time, first time.Time, last time.Time) int {
	//assume only year, month and day are set
	it := first
	i := 0
	for {
		if it.Equal(given) {
			return i
		}
		if it.Equal(last) {
			break
		}
		i++
		// incrementDay(&it)
		it = it.AddDate(0, 0, 1)
	}
	return -1
}

// calendar end: placeholder value means today
func calcLast(s string) time.Time {
	if s == "-" {
		now := time.Now()
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	}
	return parseDateStamp(s)
}
