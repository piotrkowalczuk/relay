package relay

import (
	"net"
	"strings"

	"github.com/sorcix/irc"
)

// Request represents an IRC request received by a client.
// Its handy wrapper around irc.Message object.
type Request struct {
	*irc.Message
	RemoteAddr net.Addr
	LocalAddr  net.Addr
}

// Receivers returns slice of possible recipients.
// It can be list of channels or a user that sends a message.
//
// PRIVMSG and NOTICE are the only messages available
// which actually perform delivery of a text message
// from one client to another - the rest just make it possible and try
// to ensure it happens in a reliable and structured manner.
func (r *Request) Receivers() []string {
	receivers := []string{}

	if r.Command != irc.PRIVMSG && r.Command != irc.NOTICE {
		return receivers
	}

	receivers = make([]string, 0, 1)

	for _, rec := range r.Message.Params {
		if strings.HasPrefix(rec, "#") {
			receivers = append(receivers, rec)
		}
	}

	if r.Prefix != nil && len(receivers) == 0 {
		receivers = append(receivers, r.Prefix.Name)
	}

	return receivers
}
