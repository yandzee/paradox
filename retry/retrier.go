package retry

import "github.com/yandzee/paradox/clock"

type Retrier struct {
	Backoff Backoff
	Decider Decider
	Clock   clock.Clock
}

type RetryContext struct {
	Attempt   int
	LastError error
}

func New(b Backoff, d Decider) *Retrier {
	return &Retrier{
		Backoff: b,
		Decider: d,
	}
}

func (r *Retrier) Do(fn func(*RetryContext) error) error {
	ctx := &RetryContext{}

	for {
		ctx.LastError = fn(ctx)

		switch r.Decider.Decide(ctx) {
		case Fail:
			return ctx.LastError
		case Finish:
			return nil
		case Retry:
			break
		default:
			break
		}

		r.Clock.Sleep(r.Backoff.Duration(ctx.Attempt))
		ctx.Attempt += 1
	}
}
