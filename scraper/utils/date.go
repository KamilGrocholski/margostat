package utils

import (
	"time"
)

const JAVASCRIPT_ISO_STRING string = "2006-01-02T15:04:05.999Z07:00"

func FormatDateToJavascriptISOString(time time.Time) string {
	return time.UTC().Format(JAVASCRIPT_ISO_STRING)
}
