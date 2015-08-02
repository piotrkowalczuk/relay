package action

import (
	"errors"
	"strings"
)

var (
	// ErrActionToShort is returned when incoming payload is to short.
	// To build Action struct it needs to be at least 2 characters long.
	// Shortest possible command is eg. !a
	ErrActionToShort = errors.New("antagonist/action: raw command too short, should have at least 2 characters")
)

// Action ...
type Action struct {
	Prefix    string
	Method    string
	Arguments Arguments
	Params    Params
}

// NewCommand ...
func NewCommand(raw string) (*Action, error) {
	if len(raw) < 2 {
		return nil, ErrActionToShort
	}

	parts := strings.Fields(raw)

	cmd := &Action{
		Prefix: parts[0][:1],
		Method: parts[0][1:],
	}

	if len(parts) > 2 {
		for _, arg := range parts[1:] {
			cmd.Arguments = append(cmd.Arguments, arg)
		}
	}

	return cmd, nil
}
