package amqp

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/Azure/go-amqp/internal/conn"
	"github.com/Azure/go-amqp/internal/session"
)

// ConnOptions contains the optional settings for configuring an AMQP connection.
type ConnOptions struct {
	// ContainerID sets the container-id to use when opening the connection.
	//
	// A container ID will be randomly generated if this option is not used.
	ContainerID string

	// HostName sets the hostname sent in the AMQP
	// Open frame and TLS ServerName (if not otherwise set).
	HostName string

	// IdleTimeout specifies the maximum period in milliseconds between
	// receiving frames from the peer.
	//
	// Specify a value less than zero to disable idle timeout.
	//
	// Default: 1 minute.
	IdleTimeout time.Duration

	// MaxFrameSize sets the maximum frame size that
	// the connection will accept.
	//
	// Must be 512 or greater.
	//
	// Default: 512.
	MaxFrameSize uint32

	// MaxSessions sets the maximum number of channels.
	// The value must be greater than zero.
	//
	// Default: 65535.
	MaxSessions uint16

	// Properties sets an entry in the connection properties map sent to the server.
	Properties map[string]interface{}

	// SASLType contains the specified SASL authentication mechanism.
	SASLType SASLType

	// Timeout configures how long to wait for the
	// server during connection establishment.
	//
	// Once the connection has been established, IdleTimeout
	// applies. If duration is zero, no timeout will be applied.
	//
	// Default: 0.
	Timeout time.Duration

	// TLSConfig sets the tls.Config to be used during
	// TLS negotiation.
	//
	// This option is for advanced usage, in most scenarios
	// providing a URL scheme of "amqps://" is sufficient.
	TLSConfig *tls.Config

	// test hook
	dialer dialer
}

// Client is an AMQP client connection.
type Client struct {
	conn *conn.Conn
}

// Dial connects to an AMQP server.
//
// If the addr includes a scheme, it must be "amqp", "amqps", or "amqp+ssl".
// If no port is provided, 5672 will be used for "amqp" and 5671 for "amqps" or "amqp+ssl".
//
// If username and password information is not empty it's used as SASL PLAIN
// credentials, equal to passing ConnSASLPlain option.
//
// opts: pass nil to accept the default values.
func Dial(addr string, opts *ConnOptions) (*Client, error) {
	c, err := conn.Dial(addr, opts)
	if err != nil {
		return nil, err
	}
	err = c.Start()
	if err != nil {
		return nil, err
	}
	return &Client{conn: c}, nil
}

// New establishes an AMQP client connection over conn.
// opts: pass nil to accept the default values.
func New(netConn net.Conn, opts *ConnOptions) (*Client, error) {
	c, err := conn.New(netConn, opts)
	if err != nil {
		return nil, err
	}
	err = c.Start()
	if err != nil {
		return nil, err
	}
	return &Client{conn: c}, nil
}

// Close disconnects the connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

// NewSession opens a new AMQP session to the server.
// Returns ErrConnClosed if the underlying connection has been closed.
// opts: pass nil to accept the default values.
func (c *Client) NewSession(ctx context.Context, opts *SessionOptions) (*Session, error) {
	s, err := c.conn.NewSession((*session.SessionOptions)(opts))
	if err != nil {
		return nil, err
	}

	if err = s.Begin(ctx); err != nil {
		return nil, err
	}
	return &Session{impl: s}, nil
}

// SessionOption contains the optional settings for configuring an AMQP session.
type SessionOptions struct {
	// IncomingWindow sets the maximum number of unacknowledged
	// transfer frames the server can send.
	IncomingWindow uint32

	// OutgoingWindow sets the maximum number of unacknowledged
	// transfer frames the client can send.
	OutgoingWindow uint32

	// MaxLinks sets the maximum number of links (Senders/Receivers)
	// allowed on the session.
	//
	// Minimum: 1.
	// Default: 4294967295.
	MaxLinks uint32
}
