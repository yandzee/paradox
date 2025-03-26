package clock

import "time"

type Clock interface {
	Now() int64
	Sleep(time.Duration)
	SleepChannel(time.Duration) <-chan time.Time
}
