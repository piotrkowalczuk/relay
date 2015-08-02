package relay_test

import (
	"testing"

	"github.com/piotrkowalczuk/relay"
	"github.com/sorcix/irc"
	"github.com/stretchr/testify/assert"
)

func TestJoinMessage(t *testing.T) {
	success := []struct {
		channel string
		key     string
		params  []string
		command string
	}{
		{
			channel: "#test",
			key:     "",
			params:  []string{"#test"},
			command: irc.JOIN,
		},
		{
			channel: "#test",
			key:     "key",
			params:  []string{"#test,key"},
			command: irc.JOIN,
		},
	}

	for _, data := range success {
		channel := relay.NewChannel(data.channel, data.key)
		message := relay.JoinMessage(channel)

		assert.Equal(t, data.command, message.Command)
		assert.Equal(t, data.params, message.Params)
		assert.Empty(t, message.Trailing)
	}
}
