package amqp

import (
	"context"

	"github.com/Azure/go-amqp/internal/encoding"
	"github.com/Azure/go-amqp/internal/link"
)

// Receiver receives messages on a single AMQP link.
type Receiver struct {
	impl *link.Receiver
}

// IssueCredit adds credits to be requested in the next flow
// request.
func (r *Receiver) IssueCredit(credit uint32) error {
	return r.impl.IssueCredit(credit)
}

// DrainCredit sets the drain flag on the next flow frame and
// waits for the drain to be acknowledged.
func (r *Receiver) DrainCredit(ctx context.Context) error {
	return r.impl.DrainCredit(ctx)
}

// Prefetched returns the next message that is stored in the Receiver's
// prefetch cache. It does NOT wait for the remote sender to send messages
// and returns immediately if the prefetch cache is empty. To receive from the
// prefetch and wait for messages from the remote Sender use `Receive`.
//
// When using ModeSecond, you *must* take an action on the message by calling
// one of the following: AcceptMessage, RejectMessage, ReleaseMessage, ModifyMessage.
// When using ModeFirst, the message is spontaneously Accepted at reception.
func (r *Receiver) Prefetched() *Message {
	return r.impl.Prefetched()
}

// Receive returns the next message from the sender.
//
// Blocks until a message is received, ctx completes, or an error occurs.
// When using ModeSecond, you *must* take an action on the message by calling
// one of the following: AcceptMessage, RejectMessage, ReleaseMessage, ModifyMessage.
// When using ModeFirst, the message is spontaneously Accepted at reception.
func (r *Receiver) Receive(ctx context.Context) (*Message, error) {
	return r.impl.Receive(ctx)
}

// Accept notifies the server that the message has been
// accepted and does not require redelivery.
func (r *Receiver) AcceptMessage(ctx context.Context, msg *Message) error {
	if !msg.shouldSendDisposition() {
		return nil
	}
	return r.messageDisposition(ctx, msg, &encoding.StateAccepted{})
}

// Reject notifies the server that the message is invalid.
//
// Rejection error is optional.
func (r *Receiver) RejectMessage(ctx context.Context, msg *Message, e *Error) error {
	if !msg.shouldSendDisposition() {
		return nil
	}
	return r.messageDisposition(ctx, msg, &encoding.StateRejected{Error: e})
}

// Release releases the message back to the server. The message
// may be redelivered to this or another consumer.
func (r *Receiver) ReleaseMessage(ctx context.Context, msg *Message) error {
	if !msg.shouldSendDisposition() {
		return nil
	}
	return r.messageDisposition(ctx, msg, &encoding.StateReleased{})
}

// Modify notifies the server that the message was not acted upon and should be modifed.
func (r *Receiver) ModifyMessage(ctx context.Context, msg *Message, options *ModifyMessageOptions) error {
	if !msg.shouldSendDisposition() {
		return nil
	}
	if options == nil {
		options = &ModifyMessageOptions{}
	}
	return r.messageDisposition(ctx,
		msg, &encoding.StateModified{
			DeliveryFailed:     options.DeliveryFailed,
			UndeliverableHere:  options.UndeliverableHere,
			MessageAnnotations: options.Annotations,
		})
}

// ModifyMessageOptions contains the optional parameters to ModifyMessage.
type ModifyMessageOptions struct {
	// DeliveryFailed indicates that the server must consider this an
	// unsuccessful delivery attempt and increment the delivery count.
	DeliveryFailed bool

	// UndeliverableHere indicates that the server must not redeliver
	// the message to this link.
	UndeliverableHere bool

	// Annotations is an optional annotation map to be merged
	// with the existing message annotations, overwriting existing keys
	// if necessary.
	Annotations Annotations
}

// Address returns the link's address.
func (r *Receiver) Address() string {
	return r.impl.Address()
}

// LinkName returns associated link name or an empty string if link is not defined.
func (r *Receiver) LinkName() string {
	return r.impl.LinkName()
}

// LinkSourceFilterValue retrieves the specified link source filter value or nil if it doesn't exist.
func (r *Receiver) LinkSourceFilterValue(name string) interface{} {
	return r.impl.LinkSourceFilterValue(name)
}

// Close closes the Receiver and AMQP link.
//
// If ctx expires while waiting for servers response, ctx.Err() will be returned.
// The session will continue to wait for the response until the Session or Client
// is closed.
func (r *Receiver) Close(ctx context.Context) error {
	return r.impl.Close(ctx)
}
