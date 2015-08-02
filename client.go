package antagonist

import (
	"errors"
	"net"

	"github.com/sorcix/irc"
)

var (
	// ErrJoinNotEnoughChannels ...
	ErrJoinNotEnoughChannels = errors.New("antagonist: join not enough channels")
)

// ClientOpts ...
type ClientOpts struct {
	Logger              StdLogger
	UserNick            string
	UserPassword        string
	UserMode            rune
	UserRealName        string
	RequestInterceptor  func(*irc.Message) error
	ResponseInterceptor func(*irc.Message) error
}

// Client ...
type Client struct {
	userNick            string
	userPassword        string
	userMode            rune
	userRealName        string
	remoteAddr          net.Addr
	localAddr           net.Addr
	requestInterceptor  func(*irc.Message) error
	responseInterceptor func(*irc.Message) error
	fatal               chan error
	reg                 chan bool
	connection          *irc.Conn
	logger              StdLogger
}

// NewClient ...
func NewClient(conn net.Conn, options ClientOpts) *Client {
	return &Client{
		userNick:            options.UserNick,
		userPassword:        options.UserPassword,
		userMode:            options.UserMode,
		userRealName:        options.UserRealName,
		remoteAddr:          conn.RemoteAddr(),
		localAddr:           conn.LocalAddr(),
		requestInterceptor:  options.RequestInterceptor,
		responseInterceptor: options.ResponseInterceptor,
		connection:          irc.NewConn(conn),
		logger:              options.Logger,
		fatal:               make(chan error, 1),
		reg:                 make(chan bool, 1),
	}
}

// Handle ...
func (c *Client) Handle(h Handler) {
	go c.listen(h)

	if err := c.Encode(&irc.Message{
		Command: irc.PASS,
		Params:  []string{c.userPassword},
	}); err != nil {
		c.fatal <- err
	} else {
		c.logger.Print("Message PASS has been send.")
	}

	if err := c.Encode(&irc.Message{
		Command: irc.USER,
		Params:  []string{c.userNick, string(c.userMode), "*", ":" + c.userRealName},
	}); err != nil {
		c.fatal <- err
	} else {
		c.logger.Print("Message USER has been send.")
	}

	if err := c.Encode(&irc.Message{
		Command: irc.NICK,
		Params:  []string{c.userNick},
	}); err != nil {
		c.fatal <- err
	} else {
		c.logger.Print("Message NICK has been send.")
	}
}

func (c *Client) listen(h Handler) {

	for {
		message, err := c.connection.Decode()
		if c.requestInterceptor != nil {
			err := c.requestInterceptor(message)
			if err != nil {
				break
			}
		}
		if err != nil {
			c.fatal <- err
			break
		}

		switch message.Command {
		case irc.PING:
			if err := c.Encode(&irc.Message{
				Command: irc.PONG,
			}); err != nil {
				c.fatal <- err
			}
		case irc.MODE, irc.RPL_WELCOME:
			c.logger.Print("User has been registered successfully.")
			c.reg <- true
		default:
			request := &Request{
				Message:    message,
				RemoteAddr: c.remoteAddr,
				LocalAddr:  c.localAddr,
			}

			if len(message.Params) > 0 {
				request.Channel = message.Params[0]
			} else {
				request.Channel = "<unknown>"
			}

			h.ServeIRC(newMessageWriter(c.connection), request)
		}
	}
}

// Registered ...
func (c *Client) Registered() <-chan bool {
	return c.reg
}

// Err ...
func (c *Client) Err() <-chan error {
	return c.fatal
}

// Encode ...
func (c *Client) Encode(message *irc.Message) error {
	return c.connection.Encode(message)
}

// Join ...
func (c *Client) Join(channels ...string) {
	if len(channels) == 0 {
		c.fatal <- ErrJoinNotEnoughChannels
		return
	}

	message := irc.Message{
		Command: irc.JOIN,
	}

	for _, ch := range channels {
		message.Params = append(message.Params, ch)
	}

	err := c.Encode(&message)
	if err != nil {
		c.fatal <- err
	} else {
		c.logger.Print("Message JOIN has been send.")
	}
}
