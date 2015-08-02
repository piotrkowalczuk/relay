package action

import "github.com/piotrkowalczuk/relay"

type handlers struct {
	handler Handler
	params  Params
}

// ServeMux ...
type ServeMux struct {
	handlers          map[string]handlers
	notFoundHandler   relay.Handler
	badRequestHandler relay.Handler
}

// NewServeMux ...
func NewServeMux() *ServeMux {
	return &ServeMux{
		handlers: make(map[string]handlers),
	}
}

// Handle ...
func (sm *ServeMux) Handle(method string, params Params, handler Handler) {
	sm.handlers[method] = handlers{
		handler: handler,
		params:  params,
	}
}

// NotFound ...
func (sm *ServeMux) NotFound(h relay.Handler) {
	sm.notFoundHandler = h
}

// BadRequest ...
func (sm *ServeMux) BadRequest(h relay.Handler) {
	sm.badRequestHandler = h
}

// ServeIRC ...
func (sm *ServeMux) ServeIRC(mw relay.MessageWriter, r *relay.Request) {
	cmd, err := NewCommand(r.Trailing)
	if err != nil {
		if sm.badRequestHandler != nil {
			sm.badRequestHandler.ServeIRC(mw, r)
		}
	}

	h, ok := sm.handlers[cmd.Method]
	if !ok {
		if sm.notFoundHandler != nil {
			sm.notFoundHandler.ServeIRC(mw, r)
		}

		return
	}

	h.handler.ServeIRC(mw, r, cmd)
}
