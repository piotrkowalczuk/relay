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

// Receivers ...
//
// PRIVMSG and NOTICE are the only messages available
// which actually perform delivery of a text message
// from one client to another - the rest just make it possible and try
// to ensure it happens in a reliable and structured manner.
func (r *Request) Receivers() (receivers []string) {
	if r.Command != irc.PRIVMSG && r.Command != irc.NOTICE {
		return
	}

	receivers = make([]string, 0, len(receivers))
	sendBack := false
	for _, rec := range r.Message.Params {
		if strings.HasPrefix(rec, "#") {
			receivers = append(receivers, rec)
		} else {
			if sendBack == true {
				continue
			}

			sendBack = true
			receivers = append(receivers, r.Prefix.Name)
		}
	}

	return
}
