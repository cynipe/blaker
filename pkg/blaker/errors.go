package blaker

import (
	"fmt"
	"strings"
	"time"
)

type SkipError struct {
	breakTime time.Time
	input     *RunCmdInput
}

func NewSkipError(breakTime time.Time, input *RunCmdInput) *SkipError {
	return &SkipError{
		breakTime: breakTime,
		input:     input,
	}
}

func (s *SkipError) Error() string {
	return fmt.Sprintf("the command cannot be run after %s. skipped command: `%s %s`",
		s.breakTime.Format(time.RFC3339),
		s.input.Command,
		strings.Join(s.input.Args, " "),
	)
}
