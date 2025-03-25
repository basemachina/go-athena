package athena

// PollMode is the mode of polling for query results.
type PollMode int

const (
	// PollModeConstant is the mode of polling for query results in constant intervals.
	PollModeConstant PollMode = 0

	// PollModeExponential is the mode of polling for query results in exponential intervals.
	PollModeExponential PollMode = 1
)
