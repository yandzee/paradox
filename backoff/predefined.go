package backoff

import "time"

type PredefinedBackoff struct {
	Delays []time.Duration
	Jitter float64
}

func (pb *PredefinedBackoff) Duration(i int) time.Duration {
	var d time.Duration

	n := len(pb.Delays)
	if i < n {
		d = pb.Delays[i]
	} else {
		d = pb.Delays[n-1]
	}

	return d + time.Duration(Jitter(float64(d), pb.Jitter))
}
