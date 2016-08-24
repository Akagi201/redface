// Package server implements a redis server
package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/mediocregopher/radix.v2/redis"
)

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers. If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler that calls f.
type HandlerFunc func(conn net.Conn, args []string) (interface{}, error)

// Server A Server defines parameters for running an Redis API Server.
// The zero value for Server is a valid configuration.
type Server struct {
	// "unix" for unix socket, "tcp" for tcp
	Proto string
	// TCP address to listen on, ":6379" if empty
	Addr string

	mu sync.RWMutex

	// map cmd: HandlerFunc
	m map[string]HandlerFunc
}

// Opts are Options which can be passed in to NewWithOpts. If any are set to
// their zero value the default value will be used instead
type Opts struct {
	// "unix" for unix socket, "tcp" for tcp
	Proto string
	// The address to listen at
	Host string
	// The port to listen at
	Port int
}

// New creates a new tcp redis server with the port.
func New(port int) (*Server, error) {
	return NewWithOpts(Opts{
		Port: port,
	})
}

// NewWithOpts is the same as New, but with more fine-tuned
// configuration options. See Opts for more available options.
func NewWithOpts(o Opts) (*Server, error) {
	if len(o.Proto) == 0 {
		o.Proto = "tcp"
	}

	if len(o.Host) == 0 {
		o.Host = "127.0.0.1"
	}

	if o.Port == 0 {
		o.Port = 6389
	}

	srv := &Server{
		Proto: o.Proto,
		m:     make(map[string]HandlerFunc),
	}

	if srv.Proto == "unix" {
		srv.Addr = o.Host
	} else {
		srv.Addr = fmt.Sprintf("%s:%d", o.Host, o.Port)
	}

	return srv, nil
}

// ListenAndServe listens on the TCP or Unix socket network address srv.Addr and then
// calls Serve to handle requests on incoming connections.
// If srv.Addr is blank, ":6379" is used.
// ListenAndServe always returns a non-nil error.
func (srv *Server) ListenAndServe() error {
	addr := srv.Addr
	if srv.Proto == "" {
		srv.Proto = "tcp"
	}
	if srv.Proto == "unix" && addr == "" {
		addr = "/tmp/redis.sock"
	} else if addr == "" {
		addr = ":6379"
	}
	ln, err := net.Listen(srv.Proto, addr)
	if err != nil {
		return err
	}
	return srv.Serve(ln)
}

// Serve accepts incoming connections on the Listener l, creating a new service goroutine for each.
// The service goroutines read requests and then call srv.Handler to reply to them.
//
// Serve always returns a non-nil error.
func (srv *Server) Serve(l net.Listener) error {
	defer l.Close()

	for {
		rw, err := l.Accept()
		if err != nil {
			return err
		}
		go srv.doServe(rw)
	}
}

var invalidCmdResp = redis.NewResp(errors.New("ERR invalid command"))

// doServe starts a new redis session, using `conn` as a transport.
// It reads commands using the redis protocol, passes them to `handler`,
// and returns the result.
func (srv *Server) doServe(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("redis: panic doServe %v", err)
		}
		conn.Close()
	}()

	rr := redis.NewRespReader(conn)

outer:
	for {
		var cmd string
		var args []string

		m := rr.Read()
		if m.IsType(redis.IOErr) {
			log.Printf("Client connection error %q", m.Err)
			conn.Close()
			return
		}

		parts, err := m.Array()
		if err != nil {
			log.Printf("Error parsing to array: %q", err)
			continue outer
		}

		for i := range parts {
			val, err := parts[i].Str()
			if err != nil {
				log.Printf("Invalid command part %#v: %s", parts[i], err)
				invalidCmdResp.WriteTo(conn)
				continue outer
			}
			if i == 0 {
				cmd = val
			} else {
				args = append(args, val)
			}
		}

		log.Printf("redis: doServe cmd: %s, args: %#v", cmd, args)
		srv.Dispatch(conn, cmd, args)
	}
}

// Dispatch takes in a client whose command has already been read off the
// socket, a list of arguments from that command (not including the command name
// itself), and handles that command
func (srv *Server) Dispatch(conn net.Conn, cmd string, args []string) {
	handlerFunc, ok := srv.m[strings.ToUpper(cmd)]
	if !ok {
		writeErrf(conn, "ERR unknown command %q", cmd)
		return
	}

	ret, err := handlerFunc(conn, args)
	if err != nil {
		writeErrf(conn, "ERR unexpected server-side error")
		log.Printf("command %s %#v err: %s", cmd, args, err)
		return
	}

	redis.NewResp(ret).WriteTo(conn)
}

// Handle registers the handler for the given cmd.
// If a handler already exists for pattern, Handle panics.
func (srv *Server) Handle(cmd string, handlerFunc HandlerFunc) {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	if cmd == "" {
		panic("redis: invalid cmd " + cmd)
	}

	if srv.m == nil {
		srv.m = make(map[string]HandlerFunc)
	}

	srv.m[strings.ToUpper(cmd)] = handlerFunc
}

func writeErrf(w io.Writer, format string, args ...interface{}) {
	err := fmt.Errorf(format, args...)
	redis.NewResp(err).WriteTo(w)
}
