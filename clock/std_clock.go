package clock

import "time"

type StdClock struct{}

func (sc *StdClock) Now() int64 {
	return time.Now().UnixNano()
}

func (sc *StdClock) Sleep(d time.Duration) {
	time.Sleep(d)
}

func (sc *StdClock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}
