package retry

import (
	"context"
	"errors"
	"time"

	"github.com/yandzee/paradox/backoff"
	"github.com/yandzee/paradox/clock"
)

type Retrier struct {
	backoff backoff.Backoff
	decider Decider
	clock   clock.Clock
}

type DoFn func(*RetryContext) error
type DoManualFn func(*RetryContext) (time.Duration, error)

func New(b backoff.Backoff, d Decider) Retrier {
	return NewWithClock(b, d, clock.Std)
}

func NewWithClock(b backoff.Backoff, d Decider, c clock.Clock) Retrier {
	return Retrier{
		backoff: b,
		decider: d,
		clock:   c,
	}
}

func (r *Retrier) Do(ctx context.Context, fn DoFn) error {
	return r.DoManual(ctx, func(rc *RetryContext) (time.Duration, error) {
		return 0, fn(rc)
	})
}

// NOTE: Useful when at some point during retry attempts you obtain the knowledge
// about delay that should be on the next iteration. Example: you are doing
// http requests and get rate limit error, telling you how much time you should
// wait before next attempt.
func (r *Retrier) DoManual(ctx context.Context, fn DoManualFn) error {
	rctx := RetryContext{}
	dur := time.Duration(0)

	if ctx == nil {
		ctx = context.Background()
	}

	for {
		dur, rctx.LastError = fn(&rctx)

		switch r.decider.Decide(&rctx) {
		case Fail:
			return rctx.LastError
		case Success:
			return nil
		default:
		}

		// Sleep handling
		if dur <= 0 {
			dur = r.backoff.Duration(rctx.Attempt)
		}

		select {
		case <-r.clock.After(dur):
		case <-ctx.Done():
			return errors.Join(ctx.Err(), rctx.LastError)
		}

		rctx.Attempt += 1
	}
}
