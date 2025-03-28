package backoff

import (
	"math/rand/v2"
	"time"
)

type Backoff interface {
	Duration(int) time.Duration
}

// Returns base + fraction of base proportional to factor
func Jitter(base, factor float64) float64 {
	// a := rand.Float64() gives a num in [0, 1)
	// b := 2 * a - 1 gives a num in [-1, -1)
	// c := b * factor gives a num in [-factor, +factor)

	return base*(2*rand.Float64()-1)*factor
}
