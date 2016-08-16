package main

import (
	"strings"

	"github.com/gopherjs/gopherjs/js"
)
import "github.com/googollee/go-socket.io"

func main() {
	// export the main function
	js.Global.Set("NewSocketServer", js.MakeFunc(NewServer))
}

type point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

func pointFromJS(o *js.Object) point {
	var p point
	p.X = o.Get("x").Float()
	p.Y = o.Get("y").Float()
	p.Z = o.Get("z").Float()
	return p
}

type position struct {
	Machine point `json:"machine"`
	Work    point `json:"work"`
}

func positionFromJS(o *js.Object) position {
	var p position
	p.Machine = pointFromJS(o.Get("machine"))
	p.Work = pointFromJS(o.Get("work"))
	return p
}

type portLostError struct {
	CompletedCommandCount int      `json:"completed_command_count"`
	PendingCommandCount   int      `json:"pending_command_count"`
	CurrentPosition       position `json:"current_position"`
	LastInstruction       string   `json:"last_instruction"`
	ActiveBuffer          []string `json:"active_buffer"`
	SenderNote            string   `json:"sender_note"`
}

func portLostErrorFromJS(o *js.Object) portLostError {
	var p portLostError
	p.CompletedCommandCount = o.Get("completed_command_count").Int()
	p.PendingCommandCount = o.Get("pending_command_count").Int()
	p.CurrentPosition = positionFromJS(o.Get("current_position"))
	p.LastInstruction = o.Get("last_instruction").String()
	p.ActiveBuffer = make([]string, o.Get("active_buffer").Length())
	for i := range p.ActiveBuffer {
		p.ActiveBuffer[i] = o.Get("active_buffer").Index(i).String()
	}
	p.SenderNote = o.Get("sender_note").String()
	return p
}

type machineType struct {
	Product  string `json:"product"`
	Revision string `json:"revision"`
}

func machineTypeFromJS(o *js.Object) machineType {
	var m machineType
	m.Product = o.Get("product").String()
	m.Revision = o.Get("revision").String()
	return m
}

type runTime struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

func runTimeFromJS(o *js.Object) runTime {
	var r runTime
	r.Start = o.Get("start").Call("toISOString").String()
	r.End = o.Get("end").Call("toISOString").String()
	return r
}

type serPort struct {
	VendorID     string `json:"vendorId"`
	ProductID    string `json:"productId"`
	ComName      string `json:"comName"`
	Manufacturer string `json:"manufacturer"`
}

func serPortFromJS(o *js.Object) serPort {
	var s serPort
	s.VendorID = o.Get("vendorId").String()
	s.ProductID = o.Get("productId").String()
	s.ComName = o.Get("comName").String()
	s.Manufacturer = o.Get("manufacturer").String()
	return s
}
func serPortsFromJS(o *js.Object) []serPort {
	s := make([]serPort, o.Length())
	for i := range s {
		s[i] = serPortFromJS(o.Index(i))
	}
	return s
}

type echo struct {
	Action string `json:"action"`
	Data   string `json:"data,omitempty"`
}

func echoFromJS(o *js.Object) echo {
	var e echo
	e.Action = o.Get("action").String()
	if e.Action == "read" || e.Action == "write" {
		e.Data = o.Get("data").String()
	}
	return e
}

// Emit handles all of the "sockets.emit" calls in the websocket controller
func (s *Server) Emit(this *js.Object, args []*js.Object) interface{} {
	name := args[0].String()
	switch name {
	case "ready":
		go s.sio.BroadcastTo("/", name)
	case "resumed":
		go s.sio.BroadcastTo("/", name, args[1].String(), args[2].Float())
	case "running":
		go s.sio.BroadcastTo("/", name, args[1].String(), args[2].Float())
	//case "error": not actually dispatched from anywhere
	case "port_lost":
		go s.sio.BroadcastTo("/", name, portLostErrorFromJS(args[1]))
	case "position":
		go s.sio.BroadcastTo("/", name, positionFromJS(args[1]))
	case "state":
		go s.sio.BroadcastTo("/", name, args[1].String())
	case "run-state":
		go s.sio.BroadcastTo("/", name, args[1].String())
	case "machine-settings":
		go s.sio.BroadcastTo("/", name, args[1].String())
	case "machine-type":
		go s.sio.BroadcastTo("/", name, machineTypeFromJS(args[1]))
	case "serial-number":
		go s.sio.BroadcastTo("/", name, args[1].String())
	case "run-time":
		go s.sio.BroadcastTo("/", name, runTimeFromJS(args[1]))
	case "paused":
		go s.sio.BroadcastTo("/", name, args[1].String(), args[2].Float())
	case "release":
		go s.sio.BroadcastTo("/", name, args[1].Int())
	case "stopping":
		go s.sio.BroadcastTo("/", name)
	case "grbl-error":
		go s.sio.BroadcastTo("/", name, args[1].String())
	case "grbl-alarm":
		go s.sio.BroadcastTo("/", name, args[1].String())
	case "connection_status":
		go s.sio.BroadcastTo("/", name, args[1].String())
	case "ports":
		go s.sio.BroadcastTo("/", name, serPortsFromJS(args[1]))
	case "echo":
		go s.sio.BroadcastTo("/", name, echoFromJS(args[1]))
	default:
		panic("unknown event: " + name)
	}

	return nil
}

type gcodeParam struct {
	Name  string
	GCode string
}

func (g gcodeParam) JS() *js.Object {
	n := new(js.Object)
	n.Set("name", g.Name)
	n.Set("gcode", g.GCode)
	return n
}

type setConfig struct {
	Name            string            `json:"name"`
	GCode           map[string]string `json:"gcode"`
	Baud            int               `json:"baud"`
	Separator       string            `json:"separator"`
	ReadyResponses  []string          `json:"readyResponses"`
	SuccessResponse string            `json:"successResponse"`
}

func (s setConfig) JS() *js.Object {
	n := new(js.Object)
	n.Set("name", s.Name)
	m := new(js.Object)
	for key, val := range s.GCode {
		m.Set(key, val)
	}
	n.Set("gcode", m)
	n.Set("baud", s.Baud)
	n.Set("separator", s.Separator)

	r := js.Global.Get("Array").New()
	for i, val := range s.ReadyResponses {
		r.SetIndex(i, val)
	}
	n.Set("readyResponses", r)
	n.Set("successResponse", s.SuccessResponse)
	return n
}

// On handles all of the "sockets.on" calls in websocket controller
// really just `.on("connection"...` but it handles all of the per-socket calls too
func (s *Server) On(this *js.Object, args []*js.Object) interface{} {
	if args[0].String() != "connection" {
		panic("unexpected event binding: " + args[0].String())
	}
	cb := args[1]

	s.sio.On("connection", func(sock socketio.Socket) {
		sock.Join("/")
		sock.Join("")
		o := new(js.Object)
		o.Set("on", js.MakeFunc(func(this *js.Object, args []*js.Object) interface{} {
			name := args[0].String()
			fn := args[1]
			switch name {
			case "get_connection":
				sock.On(name, func(val string) { fn.Invoke() })
			case "get_job_status":
				sock.On(name, func(val string) { fn.Invoke() })
			case "gcode":
				sock.On(name, func(val gcodeParam) { fn.Invoke(val.JS()) })
			case "get_ports":
				sock.On(name, func(val string) { fn.Invoke() })
			case "console":
				sock.On(name, func(val string) { fn.Invoke(val) })
			case "execute":
				sock.On(name, func(val []string) { fn.Invoke(val) })
			case "state":
				sock.On(name, func(val string) { fn.Invoke() })
			case "set_config":
				sock.On(name, func(val setConfig) { fn.Invoke(val.JS()) })
			case "disconnect":
				sock.On(name, func(val string) { fn.Invoke() })
			case "init_port":
				sock.On(name, func(val string) { fn.Invoke(val) })
			case "pause":
				sock.On(name, func(val string) { fn.Invoke() })
			case "acquire":
				sock.On(name, func(val int) { fn.Invoke(val) })
			case "resume":
				sock.On(name, func(val string) { fn.Invoke() })
			case "stop":
				// stop params not used
				sock.On(name, func(val string) { fn.Invoke() })
			case "echo":
				sock.On(name, func(enabled bool) { fn.Invoke(enabled) })
			case "machine-settings":
				sock.On(name, func(val string) { fn.Invoke() })
			case "sent_feedback":
				sock.On(name, func(val string) { fn.Invoke() })
			default:
				panic("unknown event: " + name)
			}
			return nil
		}))
		o.Set("emit", js.MakeFunc(func(this *js.Object, args []*js.Object) interface{} {
			event := args[0].String()
			switch event {
			case "version":
				sock.Emit("version", args[1].String())
			case "iris-state":
				sock.Emit("iris-state", args[1].String())
			default:
				panic("unknown event: " + event)
			}
			return nil
		}))
		cb.Invoke(o)
	})
	return nil
}

// SetOrigins is only called once, but allows updating the allowed origins
func (s *Server) SetOrigins(this *js.Object, args []*js.Object) interface{} {
	o := strings.Split(args[0].String(), " ")
	s.origins = o[:0]
	for _, name := range o {
		if strings.HasSuffix(name, ":80") {
			s.origins = append(s.origins, "http://"+strings.TrimSuffix(name, ":80"))
		} else if strings.HasSuffix(name, ":443") {
			s.origins = append(s.origins, "https://"+strings.TrimSuffix(name, ":443"))
		} else {
			s.origins = append(s.origins, name)
		}
	}
	return nil
}
func log(args ...interface{}) {
	js.Global.Get("console").Call("log", args...)
}
