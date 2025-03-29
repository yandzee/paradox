package tests

import (
	"context"
	"testing"
	"time"

	"github.com/yandzee/paradox/clock"
)

// Test that at least smth is not crashing
func TestStdClock(t *testing.T) {
	c := clock.Std

	if ts := c.Now(); ts <= 0 {
		t.Fatalf("Real clock gives wrong Now()")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	select {
	case <-c.After(time.Nanosecond):
	case <-ctx.Done():
		t.Fatalf("Sleep channel doesn't fire after 1 ns")
	}
}

func TestMockFastForward(t *testing.T) {
	c := clock.MockClock{}

	chans := []<-chan time.Time{
		c.After(time.Nanosecond),
		c.After(time.Second),
		c.After(time.Hour),
		c.After(time.Hour),
	}

	for i, ch := range chans {
		select {
		case <-ch:
			t.Fatalf("Channel %d is fired by no cause", i)
		default:
		}
	}

	c.FastForward()

	for i, ch := range chans {
		select {
		case now := <-ch:
			if ts := now.UnixNano(); ts != time.Hour.Nanoseconds() {
				t.Fatalf("Channel %d has received wrong now value: %d", i, ts)
			}
		default:
			t.Fatalf("Channel %d is not fired after FastForward", i)
		}
	}
}

func TestMockAdvance(t *testing.T) {
	c := clock.MockClock{}

	secondCh := c.After(time.Second)
	hourCh := c.After(time.Hour)

	for i, ch := range []<-chan time.Time{secondCh, hourCh} {
		select {
		case <-ch:
			t.Fatalf("Channel %d is fired by no cause", i)
		default:
		}
	}

	c.Advance(time.Minute)

	select {
	case now := <-secondCh:
		if ts := now.UnixNano(); ts != time.Minute.Nanoseconds() {
			t.Fatalf("Second channel has received wrong now value: %d", ts)
		}
	default:
		t.Fatalf("Second channel is not fired after Advance for minute")
	}

	select {
	case <-hourCh:
		t.Fatalf("Hour channel is fired by no cause")
	default:
	}
}
