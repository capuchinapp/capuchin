package util

import (
	"time"
)

func TimeNowUtc() time.Time {
	return time.Now().UTC()
}
