package netutil

import (
	"crypto/tls"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gopherd/log"
)

func IP(req *http.Request) string {
	ips := getproxy(req)
	if len(ips) > 0 && ips[0] != "" {
		rip := strings.Split(ips[0], ":")
		if len(rip) > 0 {
			return rip[0]
		}
	}
	ip := strings.Split(req.RemoteAddr, ":")
	if len(ip) > 0 {
		if ip[0] != "[" {
			return ip[0]
		}
	}
	return "127.0.0.1"
}

func getproxy(req *http.Request) []string {
	if ips := req.Header.Get("X-Forwarded-For"); ips != "" {
		return strings.Split(ips, ",")
	}
	return []string{}
}

// TCPKeepAliveListener wraps TCPListener with a keepalive duration
type TCPKeepAliveListener struct {
	*net.TCPListener
	duration time.Duration
}

// NewTCPKeepAliveListener creates a TCPKeepAliveListener
func NewTCPKeepAliveListener(ln *net.TCPListener, d time.Duration) *TCPKeepAliveListener {
	return &TCPKeepAliveListener{
		TCPListener: ln,
		duration:    d,
	}
}

// Accept implements net.Listener Accept method
func (ln TCPKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	if ln.duration == 0 {
		ln.duration = 3 * time.Minute
	}
	tc.SetKeepAlivePeriod(ln.duration)
	return tc, nil
}

// KeepAliveTCPConn sets conn's keepalive duration
func KeepAliveTCPConn(conn net.Conn, d time.Duration) {
	tc, ok := conn.(*net.TCPConn)
	if ok {
		tc.SetKeepAlive(true)
		tc.SetKeepAlivePeriod(d)
	}
}

// ConnHandler handles net.Conn
type ConnHandler func(ip string, conn net.Conn)

// TCPServer represents a tcp server
type TCPServer struct {
	addr     string
	handler  ConnHandler
	listener net.Listener
}

// ListenTCP creates a tcp server
func ListenTCP(addr string, handler ConnHandler, keepalive time.Duration, certs ...tls.Certificate) (*TCPServer, net.Listener, error) {
	var (
		listener net.Listener
		err      error
	)
	if len(certs) > 0 {
		config := &tls.Config{Certificates: certs}
		listener, err = tls.Listen("tcp", addr, config)
	} else {
		var a *net.TCPAddr
		a, err = net.ResolveTCPAddr("tcp", addr)
		if err == nil {
			listener, err = net.ListenTCP("tcp", a)
		}
	}
	if err != nil {
		return nil, nil, err
	}
	if keepalive > 0 {
		if l, ok := listener.(*net.TCPListener); ok {
			listener = NewTCPKeepAliveListener(l, keepalive)
		} else {
			log.Warn().Print("TCPServer.ListenAndServe: keepalive is not supported")
		}
	}
	server := new(TCPServer)
	server.addr = addr
	server.handler = handler
	return server, listener, nil
}

func (server *TCPServer) Serve(listener net.Listener) error {
	server.listener = listener
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Info().
					Error("error", err).
					String("delay", tempDelay.String()).
					Print("accept connection error, retrying")
				time.Sleep(tempDelay)
				continue
			}
			return err
		}
		tempDelay = 0
		var ip string
		if addr, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
			ip = addr.IP.String()
		}
		go server.handler(ip, conn)
	}
	return nil
}

// Shutdown shutdowns the tcp server
func (server *TCPServer) Shutdown() error {
	return server.listener.Close()
}

// ListenAndServeTCP listen and serve a tcp address
func ListenAndServeTCP(addr string, keepalive time.Duration, handler ConnHandler, certs ...tls.Certificate) error {
	server, listener, err := ListenTCP(addr, handler, keepalive, certs...)
	if err != nil {
		return err
	}
	return server.Serve(listener)
}
