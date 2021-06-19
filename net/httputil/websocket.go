package httputil

import (
	"net"
	"net/http"
	"time"

	"golang.org/x/net/websocket"

	"github.com/gopherd/doge/net/netutil"
)

// ListenAndServeWebsocket starts a http server and registers a websocket handler
func ListenAndServeWebsocket(addr, path string, handler netutil.ConnHandler) error {
	mux := http.NewServeMux()
	mux.Handle(path, websocket.Handler(func(conn *websocket.Conn) {
		handler(IP(conn.Request()), conn)
	}))
	httpServer := &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return httpServer.ListenAndServe()
}

// ListenWebsocket creates a http server and registers a websocket handler
func ListenWebsocket(addr, path string, handler netutil.ConnHandler, keepalive time.Duration) (*http.Server, net.Listener, error) {
	if addr == "" {
		addr = ":http"
	}
	mux := http.NewServeMux()
	mux.Handle(path, websocket.Handler(func(conn *websocket.Conn) {
		handler(IP(conn.Request()), conn)
	}))
	server := &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, err
	}
	if keepalive > 0 {
		ln = netutil.NewTCPKeepAliveListener(ln.(*net.TCPListener), keepalive)
	}
	return server, ln, nil
}
