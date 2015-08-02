package relay

// Channel ...
type Channel struct {
	Name string
	Key  string
}

// NewChannel ...
func NewChannel(name, key string) *Channel {
	return &Channel{
		Name: name,
		Key:  key,
	}
}

// String implements fmt.Stringer interface.
func (c *Channel) String() string {
	if c.Key == "" {
		return c.Name
	}

	return c.Name + "," + c.Key
}
