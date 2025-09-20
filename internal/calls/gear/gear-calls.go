package gear_calls

import (
	"github.com/misnaged/gear-go/internal/calls"
)

type GearCalls struct {
	c *calls.Calls
}

func New(c *calls.Calls) IGearCalls {
	return &GearCalls{c: c}
}
