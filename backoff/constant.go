package backoff

import "time"

type ConstantBackoff struct {
	Delay  time.Duration
	Jitter float64
}

func (cb *ConstantBackoff) Duration(_ int) time.Duration {
	return cb.Delay + time.Duration(Jitter(float64(cb.Delay), cb.Jitter))
}
