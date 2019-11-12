package flip

import "context"

// A cleanup function taking a context.Context only.
type Cleanup func(context.Context)

type runCleanupFunc func(ExitStatus, context.Context) int

// An interface for post-command actions.
type Cleaner interface {
	SetCleanup(ExitStatus, ...Cleanup)
	RunCleanup(ExitStatus, context.Context) int
}

type cleaner struct {
	cfns map[ExitStatus][]Cleanup
}

func newCleaner() *cleaner {
	return &cleaner{make(map[ExitStatus][]Cleanup)}
}

// Set the provided Cleanup functions to be run on the provided ExitStatus.
func (c *cleaner) SetCleanup(e ExitStatus, cfns ...Cleanup) {
	if c.cfns[e] == nil {
		c.cfns[e] = make([]Cleanup, 0)
	}

	c.cfns[e] = append(c.cfns[e], cfns...)
}

// Given an ExitStatus and a context.Context, runs any associated Cleanup functions.
func (c *cleaner) RunCleanup(e ExitStatus, ctx context.Context) int {
	if cfns, ok := c.cfns[e]; ok {
		for _, cfn := range cfns {
			cfn(ctx)
		}
	}
	if afns, ok := c.cfns[ExitAny]; ok {
		for _, afn := range afns {
			afn(ctx)
		}
	}

	return int(e)
}
