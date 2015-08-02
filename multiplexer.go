package antagonist

// ServeMux is an IRC request multiplexer. It matches the IRC command of each incoming request
// against a list of registered and calls the handler matches the given command.
// TODO(piotr): mutexes!
type ServeMux struct {
	handlers        map[string]Handler
	notFoundHandler Handler
}

// NewServeMux allocates and returns a new ServeMux.
func NewServeMux() *ServeMux {
	return &ServeMux{
		handlers: make(map[string]Handler),
	}
}

// Handle registers the handler for the given IRC command. If a handler already exists for pattern, Handle panics.
func (sm *ServeMux) Handle(command string, handler Handler) {
	if command == "" {
		panic("antagonist: missing command")
	}
	if handler == nil {
		panic("antagonist: nil handler")
	}
	if _, exists := sm.handlers[command]; exists {
		panic("antagonist: multiple registrations for " + command)
	}

	sm.handlers[command] = handler
}

// NotFound sets configurable Handler which is called when no matching Handlers is
// found.
func (sm *ServeMux) NotFound(h Handler) {
	sm.notFoundHandler = h
}

// ServeIRC dispatches the request to the handler command matches the incoming message command.
func (sm *ServeMux) ServeIRC(ew MessageWriter, r *Request) {
	h, ok := sm.handlers[r.Message.Command]
	if !ok {
		if sm.notFoundHandler != nil {
			sm.notFoundHandler.ServeIRC(ew, r)
		}

		return
	}

	h.ServeIRC(ew, r)
}
