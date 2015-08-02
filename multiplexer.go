package relay

import (
	"errors"
	"sync"
)

var (
// errLackOfMatchingHandlers is returned by handler when there is no matching handler in a map,
	errLackOfMatchingHandlers = errors.New("antagonist: lack of matching handlers")
)

// ServeMux is an IRC request multiplexer. It matches the IRC command of each incoming request
// against a list of registered and calls the handler matches the given command.
// TODO(piotr): mutexes!
type ServeMux struct {
	mu              sync.RWMutex
	handlers        map[string]Handler
	notFoundHandler Handler
}

// NewServeMux allocates and returns a new ServeMux.
func NewServeMux() *ServeMux {
	return &ServeMux{
		handlers: make(map[string]Handler),
	}
}

// Handle registers the handler for the given IRC command.
// If a handler already exists for pattern, Handle panics.
func (sm *ServeMux) Handle(command string, handler Handler) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

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

func (sm *ServeMux) handler(command string) (Handler, error) {
	h, ok := sm.handlers[command]
	if !ok {
		if sm.notFoundHandler != nil {
			return sm.notFoundHandler, nil
		}

		return nil, errLackOfMatchingHandlers
	}

	return h, nil
}

// ServeIRC dispatches the request to the handler command matches the incoming message command.
func (sm *ServeMux) ServeIRC(ew MessageWriter, r *Request) {
	h, err := sm.handler(r.Message.Command)
	if err != nil {
		if err != errLackOfMatchingHandlers {
			panic(err)
		}

		return
	}

	h.ServeIRC(ew, r)
}
