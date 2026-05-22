package time

import (
	"time"
)

type RealTimeProvider struct{}

func (r RealTimeProvider) NowUTC() time.Time {
	return time.Now().UTC()
}
