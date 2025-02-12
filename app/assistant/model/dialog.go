package model

import "context"

// IDialog is an interface to define the methods of a dialog.
type IDialog interface {
	Send(string)
	Close()
}

// Dialog is a struct that represents a dialog. A simple implementation of IDialog.
type Dialog struct {
	IDialog

	ctx    context.Context
	cancel context.CancelFunc

	_receiver chan string
}

// NewDialog creates a new Dialog.
func NewDialog() *Dialog {
	ctx, cancel := context.WithCancel(context.Background())
	return &Dialog{
		_receiver: make(chan string),
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Send sends a message to the dialog.
func (d *Dialog) Send(message string) {
	select {
	case <-d.ctx.Done():
	case d._receiver <- message:
	}
}

// Close closes the dialog.
func (d *Dialog) Close() {
	close(d._receiver)
	d.cancel()
}

// NotifyOnClosed returns a channel that notifies when the dialog is closed.
func (d *Dialog) NotifyOnClosed() <-chan struct{} {
	return d.ctx.Done()
}

// NotifyOnMessage returns a channel that notifies when a message is received.
func (d *Dialog) NotifyOnMessage() <-chan string {
	return d._receiver
}
