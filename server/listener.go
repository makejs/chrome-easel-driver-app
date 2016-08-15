package main

import (
	"fmt"
	"net"

	"github.com/gopherjs/gopherjs/js"
)

// Listener turns an "accept" event stream and implements the net.Listener interface
type Listener struct {
	ch chan net.Conn
	h  *js.Object
}

// Accept will block until the next connection is available, or the listener is closed
func (l *Listener) Accept() (net.Conn, error) {
	c := <-l.ch
	if c == nil {
		return nil, fmt.Errorf("closed")
	}
	return c, nil
}

// Addr will return the server socket address
func (l *Listener) Addr() net.Addr {
	return &net.TCPAddr{IP: net.ParseIP(l.h.Get("localAddress").String()), Port: l.h.Get("localPort").Int()}
}

// Close will close the listener
func (l *Listener) Close() error {
	l.h.Call("close")
	close(l.ch)
	return nil
}

// NewListener will create a new Listener given an instance of the Server class from shim/http.js
func NewListener(httpShim *js.Object) net.Listener {
	ch := make(chan net.Conn, 5)
	httpShim.Call("on", "accept", js.MakeFunc(func(this *js.Object, args []*js.Object) interface{} {
		ch <- NewClientConn(args[0])
		return nil
	}))
	return &Listener{ch: ch, h: httpShim}
}
