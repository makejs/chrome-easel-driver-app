package main

import (
	"net/http"

	"github.com/googollee/go-socket.io"
	"github.com/gopherjs/gopherjs/js"
)

// Server handle's the real shimming of the socket.io & http parts
type Server struct {
	sio     *socketio.Server
	origins []string
}

// NewServer is able to be exported and will accept an instance of the shim/http.js
// Server class. It creates a new socketio server and returns a javascript object
// implementing the socket.io API
func NewServer(this *js.Object, args []*js.Object) interface{} {
	server, err := socketio.NewServer(nil)
	if err != nil {
		panic(err)
	}

	s := &Server{sio: server}
	l := NewListener(args[0])
	go http.Serve(l, s)
	o := new(js.Object)
	o.Set("origins", js.MakeFunc(s.SetOrigins))

	sck := new(js.Object)
	o.Set("sockets", sck)
	sck.Set("on", js.MakeFunc(s.On))
	sck.Set("emit", js.MakeFunc(s.Emit))
	return o
}

// ServeHTTP is the core http handler
func (s *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var safeOrigin bool
	// check that the Origin header is acceptable
	for _, o := range s.origins {
		if req.Header.Get("Origin") == o {
			safeOrigin = true
			break
		}
	}
	// if not, just send Forbidden
	if !safeOrigin {
		http.Error(rw, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	// Set CORS headers. since we already checked the origin header, we can just echo it back
	rw.Header().Set("Access-Control-Allow-Origin", req.Header.Get("Origin"))
	rw.Header().Set("Access-Control-Allow-Credentials", "true")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// pass the request on to the socket.io library
	s.sio.ServeHTTP(rw, req)
}
