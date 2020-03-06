package util

import "time"

func FormatDateTime(dt time.Time) string {
	return dt.Format("2006-01-02 15:04:05")
}