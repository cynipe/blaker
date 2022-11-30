package blaker

import (
	"fmt"
	"time"
)

type BreakTimeAfterError struct {
	breakTime time.Time
}

func NewBreakTimeAfterError(breakTime time.Time) *BreakTimeAfterError {
	return &BreakTimeAfterError{
		breakTime: breakTime,
	}
}

func (s *BreakTimeAfterError) Error() string {
	return fmt.Sprintf("the command cannot be run after %s.",
		s.breakTime.Format(time.RFC3339),
	)
}
