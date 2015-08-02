package action

import "github.com/piotrkowalczuk/antagonist"

type handlers struct {
	handler Handler
	params  Params
}

// ServeMux ...
type ServeMux struct {
	handlers          map[string]handlers
	notFoundHandler   antagonist.Handler
	badRequestHandler antagonist.Handler
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
func (sm *ServeMux) NotFound(h antagonist.Handler) {
	sm.notFoundHandler = h
}

// BadRequest ...
func (sm *ServeMux) BadRequest(h antagonist.Handler) {
	sm.badRequestHandler = h
}

// ServeIRC ...
func (sm *ServeMux) ServeIRC(mw antagonist.MessageWriter, r *antagonist.Request) {
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
