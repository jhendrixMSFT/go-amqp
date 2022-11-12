package amqp

import (
	"context"
	"time"

	"github.com/Azure/go-amqp/internal/session"
)

// Default session options
const (
	defaultWindow = 5000
)

// Default link options
const (
	defaultLinkCredit      = 1
	defaultLinkBatching    = false
	defaultLinkBatchMaxAge = 5 * time.Second
)

// Session is an AMQP session.
//
// A session multiplexes Receivers.
type Session struct {
	impl *session.Session
}

// Close gracefully closes the session.
//
// If ctx expires while waiting for servers response, ctx.Err() will be returned.
// The session will continue to wait for the response until the Client is closed.
func (s *Session) Close(ctx context.Context) error {
	return s.impl.Close(ctx)
}

// NewReceiver opens a new receiver link on the session.
// opts: pass nil to accept the default values.
func (s *Session) NewReceiver(ctx context.Context, source string, opts *ReceiverOptions) (*Receiver, error) {
	rcv, err := s.impl.NewReceiver(ctx, source, opts)
	if err != nil {
		return nil, err
	}
	return &Receiver{impl: rcv}, nil
}

// NewSender opens a new sender link on the session.
// opts: pass nil to accept the default values.
func (s *Session) NewSender(ctx context.Context, target string, opts *SenderOptions) (*Sender, error) {
	snd, err := s.impl.NewSender(ctx, target, opts)
	if err != nil {
		return nil, err
	}
	return &Sender{impl: snd}, nil
}
