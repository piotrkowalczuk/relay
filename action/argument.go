package action

import "strings"

// Arguments ...
type Arguments []string

// Join is a handy wrapper for strings.Join function.
func (a *Arguments) Join(separator string) string {
	return strings.Join(*a, separator)
}

// Sentence returns arguments connected with space.
func (a *Arguments) Sentence() string {
	return a.Join(" ")
}
