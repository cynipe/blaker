package blaker

import (
	"fmt"
	"strings"
	"time"
)

type BreakError struct {
	breakTime time.Time
	input     *RunCmdInput
}

func NewSkipError(breakTime time.Time, input *RunCmdInput) *BreakError {
	return &BreakError{
		breakTime: breakTime,
		input:     input,
	}
}

func (s *BreakError) Error() string {
	return fmt.Sprintf("the command cannot be run after %s. skipped command: `%s %s`",
		s.breakTime.Format(time.RFC3339),
		s.input.Command,
		strings.Join(s.input.Args, " "),
	)
}
