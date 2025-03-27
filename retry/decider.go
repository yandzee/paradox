package retry

type Decision int

const (
	Retry Decision = iota
	Fail
	Success
)

const InfiniteAttempts = 0

var DefaultDecider = &CustomDecider{
	Attempts: 5,
}

type Decider interface {
	Decide(c *RetryContext) Decision
}

type CustomDecider struct {
	Attempts    int
	DecideError func(error) Decision
}

func (dd *CustomDecider) Decide(rctx *RetryContext) Decision {
	switch {
	case rctx.LastError == nil:
		return Success
	case dd.Attempts != InfiniteAttempts && rctx.Attempt >= dd.Attempts:
		return Fail
	case dd.DecideError != nil:
		return dd.DecideError(rctx.LastError)
	default:
		return Retry
	}
}
