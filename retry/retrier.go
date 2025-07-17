package retry

import (
	"context"
	"errors"

	"github.com/yandzee/paradox/backoff"
	"github.com/yandzee/paradox/clock"
)

type Retrier struct {
	Backoff backoff.Backoff
	Decider Decider
	Clock   clock.Clock
}

type RetryContext struct {
	Attempt   int
	LastError error
}

func New(b backoff.Backoff, d Decider) *Retrier {
	return NewClock(b, d, clock.Std)
}

func NewClock(b backoff.Backoff, d Decider, c clock.Clock) *Retrier {
	return &Retrier{
		Backoff: b,
		Decider: d,
		Clock:   c,
	}
}

func (r *Retrier) Do(ctx context.Context, fn func(*RetryContext) error) error {
	rctx := &RetryContext{}

	for {
		rctx.LastError = fn(rctx)

		switch r.Decider.Decide(rctx) {
		case Fail:
			return rctx.LastError
		case Success:
			return nil
		default:
		}

		// Sleep handling
		ch := r.Clock.After(r.Backoff.Duration(rctx.Attempt))
		if ctx != nil {
			select {
			case <-ch:
			case <-ctx.Done():
				return errors.Join(ctx.Err(), rctx.LastError)
			}
		} else {
			<-ch
		}

		rctx.Attempt += 1
	}
}
