package backoff

import (
	"math"
	"time"
)

type ExponentialBackoff struct {
	Min    time.Duration
	Max    time.Duration
	Factor float64
	Jitter float64
}

func (eb *ExponentialBackoff) Duration(p int) time.Duration {
	minf, maxf := float64(eb.Min), float64(eb.Max)

	dur := math.Min(maxf, minf*math.Pow(eb.Factor, float64(p)))
	dur = Jitter(dur, eb.Jitter)

	return time.Duration(dur)
}
