package freenode

import (
	"github.com/piotrkowalczuk/relay"
	"github.com/sorcix/irc"
)

const (
	// Addr ...
	Addr = "chat.freenode.net:6697"
)

// Registrant ...
type Registrant struct {
	logger relay.StdLogger
}

// NewRegistrant ...
func NewRegistrant(logger relay.StdLogger) *Registrant {
	return &Registrant{
		logger: logger,
	}
}

// Register ...
func (r *Registrant) Register(c *relay.Client) error {
	pass := &irc.Message{
		Command: irc.PASS,
		Params:  []string{c.User.Password},
	}
	user := &irc.Message{
		Command: irc.USER,
		Params:  []string{c.User.Nick, string(c.User.Mode), "*", ":" + c.User.RealName},
	}
	nick := &irc.Message{
		Command: irc.NICK,
		Params:  []string{c.User.Nick},
	}

	if err := c.Encode(pass); err != nil {
		return err
	}
	r.logger.Print("Message PASS has been send.")

	if err := c.Encode(user); err != nil {
		return err
	}
	r.logger.Print("Message USER has been send.")

	if err := c.Encode(nick); err != nil {
		return err
	}
	r.logger.Print("Message NICK has been send.")

	return nil
}
