package retry

type Decision int

const (
	Retry Decision = iota
	Fail
	Success
)

func (d Decision) String() string {
	switch d {
	case Retry:
		return "Retry"
	case Fail:
		return "Fail"
	case Success:
		return "Success"
	}

	return "Unknown"
}
