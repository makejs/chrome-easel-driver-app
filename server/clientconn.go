package main

import (
	"io"
	"net"
	"time"

	"github.com/gopherjs/gopherjs/js"
)

// ClientConn is a tcp connection. It implements the net.Conn interface so that
// it may be used with the Go runtime net/http stack
type ClientConn struct {
	*io.PipeReader
	c          *js.Object
	data       chan []byte
	localAddr  *net.TCPAddr
	remoteAddr *net.TCPAddr
	closed     bool
}

// NewClientConn will create a new ClientConn from a javascript object
// it is expected to be an instance of the Socket class (from shims/http)
func NewClientConn(c *js.Object) *ClientConn {
	r, w := io.Pipe()
	// create a buffer for data packets
	d := make(chan []byte, 1000)
	go func() {
		// in the background we write the data from the queue
		// into the io pipe. This will exit when the channel closes
		for data := range d {
			_, err := w.Write(data)
			if err != nil {
				r.CloseWithError(err)
				c.Call("close")
				break
			}
		}
	}()
	// bind a function to the "data" event from the Socket
	c.Call("on", "data", js.MakeFunc(func(this *js.Object, args []*js.Object) interface{} {
		// make a Uint8Array view of the ArrayBuffer, which will give us a []byte
		// then add it to the queue
		data := js.Global.Get("Uint8Array").New(args[0]).Interface().([]byte)
		d <- data
		return nil
	}))
	return &ClientConn{
		PipeReader: r,
		c:          c,
		data:       d,
		localAddr:  &net.TCPAddr{IP: net.ParseIP(c.Get("localAddress").String()), Port: c.Get("localPort").Int()},
		remoteAddr: &net.TCPAddr{IP: net.ParseIP(c.Get("peerAddress").String()), Port: c.Get("peerPort").Int()},
	}
}

// Close will close the connection
func (c *ClientConn) Close() error {
	if c.closed {
		return nil
	}
	c.closed = true
	c.PipeReader.Close()
	c.c.Call("close")
	// don't forget to close the channel, so the background goroutine exits!
	close(c.data)
	return nil
}

// Write will write data to the underlying socket
func (c *ClientConn) Write(p []byte) (int, error) {
	// simply a passthrough, after turning p into an ArrayBuffer
	c.c.Call("write", js.NewArrayBuffer(p))
	return len(p), nil
}

// TODO implement SetDeadline
func (c *ClientConn) SetDeadline(t time.Time) error {
	return nil
}

// TODO implement SetReadDeadline
func (c *ClientConn) SetReadDeadline(t time.Time) error {
	return nil
}

// TODO implement SetReadDeadline
func (c *ClientConn) SetWriteDeadline(t time.Time) error {
	return nil
}

// LocalAddr will return the local TCP endpoint
func (c *ClientConn) LocalAddr() net.Addr {
	return c.localAddr
}

// RemoteAddr will return the remote TCP endpoint
func (c *ClientConn) RemoteAddr() net.Addr {
	return c.remoteAddr
}
