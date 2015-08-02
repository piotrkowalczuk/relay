package relay

import (
	"io"

	"github.com/sorcix/irc"
)

// MessageWriter interface is used by a Handler to construct an IRC message.
type MessageWriter interface {
	// Params returns the param slice that will be sent.
	// Changing the params after a call to Write has no effect.
	Params() *Params
	// Write writes the data to the message trailing as part of an IRC reply.
	// If WriteCommand has not yet been called, Write calls WriteCommand(irc.PRIVMSG)
	// before writing the data.
	Write([]byte) (int, error)
	// WriteCommand sets IRC message command type, irc.PRIVMSG by default.
	WriteCommand(string)
	WriteParams(...string)
}

type messageWriter struct {
	params  *Params
	command string
	writer  io.Writer
}

func newMessageWriter(w io.Writer) *messageWriter {
	return &messageWriter{
		params:  &Params{},
		writer:  w,
		command: irc.PRIVMSG,
	}
}

func (mw *messageWriter) Write(b []byte) (int, error) {
	m := &irc.Message{
		Params:  *mw.params,
		Command: mw.command,
	}

	if len(b) == 0 {
		m.EmptyTrailing = true
	} else {
		m.Trailing = string(b)
	}

	return mw.writer.Write(m.Bytes())
}

func (mw *messageWriter) WriteCommand(command string) {
	mw.command = command
}

func (mw *messageWriter) WriteParams(params ...string) {
	ps := make(Params, 0, len(params))
	for _, p := range params {
		ps.Set(p)
	}

	mw.params = &ps
}

func (mw *messageWriter) Params() *Params {
	return mw.params
}
