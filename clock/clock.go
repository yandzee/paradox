package clock

import "time"

var Std = &StdClock{}

type Clock interface {
	Now() int64
	Sleep(time.Duration)
	After(time.Duration) <-chan time.Time
}
