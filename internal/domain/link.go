package domain

import "time"

type Link struct {
	ID        int64     // snowflake
	Short     string    // code for short link
	LongURL   string    // original url
	CreatedAt time.Time // creation time
}
