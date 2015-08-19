package relay_test

import (
	"testing"

	"github.com/piotrkowalczuk/relay"
	"github.com/sorcix/irc"
	"github.com/stretchr/testify/assert"
)

func TestRequestReceivers(t *testing.T) {
	success := []struct {
		msg *irc.Message
		rec []string
	}{
		{
			msg: irc.ParseMessage(":Angel PRIVMSG user1, user2 :example message"),
			rec: []string{"Angel"},
		},
		{
			msg: irc.ParseMessage(":Angel PRIVMSG #general :example message"),
			rec: []string{"#general"},
		},
		{
			msg: irc.ParseMessage(":Angel PRIVMSG #general , #privatechannel :example message"),
			rec: []string{"#general", "#privatechannel"},
		},
	}

	for _, args := range success {
		req := &relay.Request{Message: args.msg}
		rec := req.Receivers()

		assert.Equal(t, args.rec, rec)
	}
}
