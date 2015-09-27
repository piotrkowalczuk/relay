package relay

import (
	"errors"
	"net"

	"strconv"

	"sync"

	"github.com/sorcix/irc"
)

const (
	successfulRegistration = "successfull registration"
)

var (
	// ErrJoinNotEnoughChannels is returned by Join if it
	// will be called without single channel as an argument.
	ErrJoinNotEnoughChannels = errors.New("relay: join not enough channels")
)

// User ...
type User struct {
	Nick     string
	Password string
	Mode     rune
	RealName string
}

// ClientOpts represent optional values that can be pass to the client.
type ClientOpts struct {
	Logger              StdLogger
	RequestInterceptor  func(*irc.Message) error
	ResponseInterceptor func(*irc.Message) error
}

// Client represents connection and logic that corresponds to single IRC server.
// It can be used to make requests.
type Client struct {
	User                *User
	remoteAddr          net.Addr
	localAddr           net.Addr
	requestInterceptor  func(*irc.Message) error
	responseInterceptor func(*irc.Message) error
	err                 chan error
	reg                 chan bool
	isRegisteredLock    sync.RWMutex
	isRegistered        bool
	connection          *irc.Conn
	logger              StdLogger
	handler             Handler
}

// NewClient ...
func NewClient(conn net.Conn, user *User) *Client {
	return NewClientWithOpts(conn, user, nil)
}

// NewClientWithOpts ...
func NewClientWithOpts(conn net.Conn, user *User, options *ClientOpts) *Client {
	c := &Client{
		User:       user,
		remoteAddr: conn.RemoteAddr(),
		localAddr:  conn.LocalAddr(),
		connection: irc.NewConn(conn),
		err:        make(chan error, 1000000),
		reg:        make(chan bool, 1),
	}

	if options != nil {
		c.requestInterceptor = options.RequestInterceptor
		c.responseInterceptor = options.ResponseInterceptor
		c.logger = options.Logger
	}

	return c
}

// ListenAndReply listens for IRC messages on TCP connection
// and handle them with provided handler. If no handler is passed it panics.
func (c *Client) ListenAndReply() {
	if c.handler == nil {
		panic("relay: handler is nil")
	}

	go c.listen(c.handler)
}

// Handle sets handler for the client.
func (c *Client) Handle(h Handler) {
	c.handler = h
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
			c.err <- err
			break
		}

		logIncomingMessage(c.logger, message)

		switch message.Command {
		case irc.PING:
			pong := &irc.Message{
				Command:  irc.PONG,
				Trailing: message.Trailing,
			}
			if err := c.Encode(pong); err != nil {
				c.err <- err
			}
		case irc.MODE, irc.RPL_WELCOME:
			if !c.IsRegistered() {
				c.isRegisteredLock.Lock()
				c.isRegistered = true
				c.isRegisteredLock.Unlock()

				log(c.logger, successfulRegistration)
				c.reg <- true
			}
		default:
			if IsErrorCommand(message.Command) {
				c.err <- &ErrorMessage{
					Message: message,
				}

				continue
			}

			if c.IsRegistered() {
				request := &Request{
					Message:    message,
					RemoteAddr: c.remoteAddr,
					LocalAddr:  c.localAddr,
				}

				h.ServeIRC(newMessageWriter(c.connection), request)
			}
		}
	}
}

// Registered ...
func (c *Client) Registered() <-chan bool {
	return c.reg
}

// IsRegistered ...
func (c *Client) IsRegistered() bool {
	c.isRegisteredLock.RLock()
	defer c.isRegisteredLock.RUnlock()

	return c.isRegistered
}

// Err returns channel
func (c *Client) Err() <-chan error {
	return c.err
}

// Encode writes the IRC encoding of irc.Message to the stream.
// Its wrapper for irc.Encoder.Encode method.
func (c *Client) Encode(message *irc.Message) error {
	logOutgoingMessage(c.logger, message)

	return c.connection.Encode(message)
}

// Join sends JOIN message to the IRC server with given channels.
func (c *Client) Join(channels ...*Channel) error {
	if len(channels) == 0 {
		return ErrJoinNotEnoughChannels
	}

	return c.Encode(JoinMessage(channels...))
}

// IsErrorCommand checks if given command reports an error.
// TODO: replace with map
func IsErrorCommand(c string) bool {
	cmd, err := strconv.ParseInt(c, 10, 64)
	if err == nil {
		return (cmd > 400 && cmd < 503) || (cmd > 903 && cmd < 908)
	}

	return false
}
