/*
Copyright 2024 The west2-online Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package model

import (
	"context"
	"sync"
)

// IDialog is an interface to define the methods of a dialog.
type IDialog interface {
	// Unique returns a unique string to identify the dialog.
	Unique() string

	// Message returns the message raising the dialog.
	Message() string

	// Send sends a message to the dialog.
	Send(msg string)

	// Close closes the dialog.
	Close()
}

// Dialog is a struct that represents a dialog. A simple implementation of IDialog.
type Dialog struct {
	IDialog

	unique  string
	message string

	ctx    context.Context
	cancel context.CancelFunc

	_receiver chan string
	_mutex    *sync.RWMutex
	_done     *bool
}

// NewDialog creates a new Dialog.
func NewDialog(id string, input string) *Dialog {
	ctx, cancel := context.WithCancel(context.Background())
	mu := &sync.RWMutex{}
	_done := false
	done := &_done
	d := &Dialog{
		ctx:       ctx,
		cancel:    cancel,
		unique:    id,
		message:   input,
		_receiver: make(chan string),
		_mutex:    mu,
		_done:     done,
	}
	go func() {
		<-ctx.Done()
		mu.Lock()
		*done = true
		mu.Unlock()
	}()
	return d
}

// Unique returns a unique string to identify the dialog.
func (d *Dialog) Unique() string {
	return d.unique
}

// Message returns the message raising the dialog.
func (d *Dialog) Message() string {
	return d.message
}

// Send sends a message to the dialog, if the dialog is closed, it will return immediately.
func (d *Dialog) Send(message string) {
	d._mutex.RLock()
	if *d._done {
		d._mutex.RUnlock()
		return
	}
	d._mutex.RUnlock()
	d._receiver <- message
}

// Close closes the dialog.
func (d *Dialog) Close() {
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
