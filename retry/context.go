package retry

import "time"

type RetryContext struct {
	Attempt      int
	LastError    error
	OneshotDelay time.Duration
}
