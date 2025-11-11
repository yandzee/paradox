# Paradox [![Go Reference](https://pkg.go.dev/badge/github.com/yandzee/paradox.svg)](https://pkg.go.dev/github.com/yandzee/paradox)

A set of time related utilities.

## Clock

This packages provides a missing `Clock` [interface](https://github.com/yandzee/paradox/blob/main/clock/clock.go)
that can be used to decouple program from a real clock functions and thus make
entire code testable.

See [this example](https://github.com/yandzee/paradox/blob/51b680a65ea80029253879550ea8c1dcfaadb60c/tests/retry_test.go#L93-L97)
of its real world usage.


## Retrier

Provides extensible `Retrier` entity, allowing to improve resilience in a program.
What makes this implementation special is its simple interface and [extensibility](https://github.com/yandzee/paradox/blob/51b680a65ea80029253879550ea8c1dcfaadb60c/retry/retrier.go#L21-L31) in terms of
[backoff](https://github.com/yandzee/paradox/blob/51b680a65ea80029253879550ea8c1dcfaadb60c/backoff/backoff.go),
[decider](https://github.com/yandzee/paradox/blob/51b680a65ea80029253879550ea8c1dcfaadb60c/retry/decider.go#L5-L7)
and [clock](https://github.com/yandzee/paradox/blob/51b680a65ea80029253879550ea8c1dcfaadb60c/clock/clock.go)
abstractions that are used under the hood.

#### Out of order delay

Another special feature is [the ability](https://github.com/yandzee/paradox/blob/51b680a65ea80029253879550ea8c1dcfaadb60c/retry/retrier.go#L43)
to run retriable code with custom logic of picking delay time between retries.
It may seem like a pretty rare case, but if you got into that, this `DoManual` function
can make your handling code much more cleaner.
