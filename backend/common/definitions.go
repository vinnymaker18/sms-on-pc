package common

import "time"

// SMSMessage represents a single SMS message in the system.
type SMSMessage = struct {
	UserID        int64
	Time          time.Time
	OriginAddress string
	MsgBody       string
}
