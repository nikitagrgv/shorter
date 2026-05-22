package domain

import "time"

type TimeProvider interface {
	NowUTC() time.Time
}
