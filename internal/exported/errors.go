package exported

import (
	"errors"
	"fmt"

	"github.com/Azure/go-amqp/internal/encoding"
)

// Errors
var (
	// ErrSessionClosed is propagated to Sender/Receivers
	// when Session.Close() is called.
	ErrSessionClosed = errors.New("amqp: session closed")

	// ErrLinkClosed is returned by send and receive operations when
	// Sender.Close() or Receiver.Close() are called.
	ErrLinkClosed = errors.New("amqp: link closed")
)

// DetachError is returned by a link (Receiver/Sender) when a detach frame is received.
//
// RemoteError will be nil if the link was detached gracefully.
type DetachError struct {
	RemoteError *encoding.Error
}

func (e *DetachError) Error() string {
	return fmt.Sprintf("link detached, reason: %+v", e.RemoteError)
}

// ConnectionError is propagated to Session and Senders/Receivers
// when the connection has been closed or is no longer functional.
type ConnectionError struct {
	inner error
}

func NewConnectionError(inner error) error {
	return &ConnectionError{inner: inner}
}

func (c *ConnectionError) Error() string {
	if c.inner == nil {
		return "amqp: connection closed"
	}
	return c.inner.Error()
}
