package clock

import (
	"sync"
	"time"
)

type MockClock struct {
	sync.RWMutex

	Current int64
	waiters waitersMap
}

type waitersMap map[int64][]chan time.Time

func (mc *MockClock) Now() (t int64) {
	mc.RLock()
	t = mc.Current
	mc.RUnlock()

	return t
}

func (mc *MockClock) Sleep(d time.Duration) {
	ch := mc.SleepChannel(d)
	<-ch
}

func (mc *MockClock) SleepChannel(d time.Duration) <-chan time.Time {
	mc.Lock()

	waiter := make(chan time.Time, 1)
	targetTime := mc.Current + d.Nanoseconds()

	if mc.waiters == nil {
		mc.waiters = make(waitersMap)
	}

	mc.waiters[targetTime] = append(mc.waiters[targetTime], waiter)

	mc.Unlock()
	return waiter
}

func (mc *MockClock) Advance(d int64) int64 {
	mc.Lock()
	mc.Current += d
	now := time.Unix(0, mc.Current)

	for targetTime, waiters := range mc.waiters {
		if targetTime > mc.Current {
			continue
		}

		for _, waiter := range waiters {
			waiter <- now
		}

		delete(mc.waiters, targetTime)
	}

	mc.Unlock()
	return mc.Current
}

func (mc *MockClock) FastForward() int64 {
	mc.Lock()
	maxTime := mc.Current

	for targetTime := range mc.waiters {
		if targetTime > maxTime {
			maxTime = targetTime
		}
	}

	mc.Current = maxTime
	now := time.Unix(0, maxTime)

	for targetTime, waiters := range mc.waiters {
		for _, waiter := range waiters {
			waiter <- now
		}

		delete(mc.waiters, targetTime)
	}

	mc.Unlock()
	return mc.Current
}
