package util

import (
	timepkg "time"
)

func Timestamp(asString bool) interface{} {
	time := timepkg.Now()
	if asString {
		return time.Format(timepkg.RFC3339Nano)
	} else {
		return time
	}
}
