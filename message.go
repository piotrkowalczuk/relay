package relay

import "github.com/sorcix/irc"

// JoinMessage builds IRC message in a way that it satisfy IRC JOIN command structure.
func JoinMessage(channels ...*Channel) *irc.Message {
	m := &irc.Message{
		Command: irc.JOIN,
		Params:  make([]string, 0, len(channels)),
	}

	for _, ch := range channels {
		m.Params = append(m.Params, ch.String())
	}

	return m
}
