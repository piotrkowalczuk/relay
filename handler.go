package antagonist

// Handler is a basic wrapper for ServeIRC method.
type Handler interface {
	ServeIRC(MessageWriter, *Request)
}
