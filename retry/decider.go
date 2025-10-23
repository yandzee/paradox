package retry

const InfiniteAttempts = 0

type Decider interface {
	Decide(c *RetryContext) Decision
}

type DecideFn func(*RetryContext) Decision

type CustomDecider struct {
	DecideFn DecideFn
}

func (dd *CustomDecider) Decide(rctx *RetryContext) Decision {
	return dd.DecideFn(rctx)
}

func NewDecider(fn DecideFn) Decider {
	return &CustomDecider{
		DecideFn: fn,
	}
}

func NewAttemptsDecider(maxAttemptIndex int) Decider {
	return NewDecider(func(rctx *RetryContext) Decision {
		switch {
		case rctx.LastError == nil:
			return Success
		case maxAttemptIndex != InfiniteAttempts && rctx.Attempt >= maxAttemptIndex:
			return Fail
		default:
			return Retry
		}
	})
}

func NewAttemptsErrDecider(maxAttemptIndex int, errFn func(error) Decision) Decider {
	return NewDecider(func(rctx *RetryContext) Decision {
		switch {
		case rctx.LastError == nil:
			return Success
		case maxAttemptIndex != InfiniteAttempts && rctx.Attempt >= maxAttemptIndex:
			return Fail
		case errFn != nil:
			return errFn(rctx.LastError)
		default:
			return Retry
		}
	})
}

func NewErrDecider(errFn func(error) Decision) Decider {
	return NewDecider(func(rctx *RetryContext) Decision {
		switch {
		case rctx.LastError == nil:
			return Success
		case errFn != nil:
			return errFn(rctx.LastError)
		default:
			return Retry
		}
	})
}
