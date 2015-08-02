package relay

import "github.com/sorcix/irc"

// ErrorMessage ...
type ErrorMessage struct {
	*irc.Message
}

// Error implements error interface.
func (em *ErrorMessage) Error() string {
	return em.String()
}
