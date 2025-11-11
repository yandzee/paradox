package clock

import (
	"sync"
	"time"
)

type MockClock struct {
	Current int64

	// NOTE: If set, every waiter creation advances current time to its desired value
	AutoFastForward bool

	mx      sync.RWMutex
	waiters waitersMap
}

type waitersMap map[int64][]chan time.Time

func (mc *MockClock) Now() (t int64) {
	mc.mx.RLock()
	t = mc.Current
	mc.mx.RUnlock()

	return t
}

func (mc *MockClock) Sleep(d time.Duration) {
	ch := mc.After(d)
	<-ch
}

func (mc *MockClock) After(d time.Duration) <-chan time.Time {
	mc.mx.Lock()

	waiter := make(chan time.Time, 1)
	targetTime := mc.Current + d.Nanoseconds()

	mc.handleWaiter(waiter, targetTime)

	mc.mx.Unlock()
	return waiter
}

func (mc *MockClock) Advance(d time.Duration) int64 {
	mc.mx.Lock()
	mc.Current += d.Nanoseconds()
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

	mc.mx.Unlock()
	return mc.Current
}

func (mc *MockClock) FastForward() int64 {
	mc.mx.Lock()
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

	mc.mx.Unlock()
	return mc.Current
}

func (mc *MockClock) handleWaiter(waiter chan time.Time, targetTime int64) {
	if mc.waiters == nil {
		mc.waiters = make(waitersMap)
	}

	if !mc.AutoFastForward {
		mc.waiters[targetTime] = append(mc.waiters[targetTime], waiter)
		return
	}

	if targetTime > mc.Current {
		mc.Current = targetTime
	}

	waiter <- time.Unix(0, mc.Current)
}
