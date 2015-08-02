package antagonist

import (
	"net"

	"github.com/sorcix/irc"
)

// Request represents an IRC request received by a server.
// Its handy wrapper around irc.Message object.
type Request struct {
	*irc.Message
	RemoteAddr net.Addr
	LocalAddr  net.Addr
	Channel    string
}
