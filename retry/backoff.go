package retry

import "time"

type Backoff interface {
	Duration(int) time.Duration
}
