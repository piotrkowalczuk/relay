package action

import "github.com/piotrkowalczuk/antagonist"

// Handler ...
type Handler interface {
	ServeIRC(antagonist.MessageWriter, *antagonist.Request, *Action)
}
