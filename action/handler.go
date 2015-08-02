package action

import (
	"fmt"

	"github.com/piotrkowalczuk/relay"
	"github.com/sorcix/irc"
)

// Handler is a simple wrapper for ServeIRC method.
// In compare to relay.Handler it pass also Action object.
type Handler interface {
	ServeIRC(relay.MessageWriter, *relay.Request, *Action)
}

// NotFound replies to the channel with an INFO command.
func NotFound(mw relay.MessageWriter, r *relay.Request, a *Action) {
	mw.WriteCommand(irc.INFO)
	mw.Params().Set(r.Receivers()...)
	fmt.Fprintf(mw, "Oops... action %s not found.", a.Method)
}

// NotFoundHandler returns a simple request handler
// that replies to each request with a ``404 page not found'' reply.
func NotFoundHandler() Handler { return HandlerFunc(NotFound) }

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as IRC handlers.  If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler object that calls f.
type HandlerFunc func(relay.MessageWriter, *relay.Request, *Action)

// ServeIRC calls f(mw, r).
func (f HandlerFunc) ServeIRC(mw relay.MessageWriter, r *relay.Request, a *Action) {
	f(mw, r, a)
}
