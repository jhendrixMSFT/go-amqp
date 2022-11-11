package amqp

import (
	"context"

	"github.com/Azure/go-amqp/internal/link"
)

// Sender sends messages on a single AMQP link.
type Sender struct {
	impl *link.Sender
}

// LinkName() is the name of the link used for this Sender.
func (s *Sender) LinkName() string {
	return s.impl.LinkName()
}

// MaxMessageSize is the maximum size of a single message.
func (s *Sender) MaxMessageSize() uint64 {
	return s.impl.MaxMessageSize()
}

// Send sends a Message.
//
// Blocks until the message is sent, ctx completes, or an error occurs.
//
// Send is safe for concurrent use. Since only a single message can be
// sent on a link at a time, this is most useful when settlement confirmation
// has been requested (receiver settle mode is "Second"). In this case,
// additional messages can be sent while the current goroutine is waiting
// for the confirmation.
func (s *Sender) Send(ctx context.Context, msg *Message) error {
	return s.impl.Send(ctx, msg)
}

// Address returns the link's address.
func (s *Sender) Address() string {
	return s.impl.Address()
}

// Close closes the Sender and AMQP link.
func (s *Sender) Close(ctx context.Context) error {
	return s.impl.Close(ctx)
}
