package retry

type Decision int

const (
	Retry Decision = iota
	Fail
	Finish
)

type Decider interface {
	Decide(c *RetryContext) Decision
}
