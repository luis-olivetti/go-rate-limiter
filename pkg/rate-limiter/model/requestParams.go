package model

import "time"

type RequestParams struct {
	Key        string
	LimitCount int64
	Interval   time.Duration
	BlockTime  time.Duration
}
