package httputil

import (
	"net"
	"net/http"
	"time"

	"golang.org/x/net/websocket"

	"github.com/gopherd/doge/net/netutil"
)

// ListenAndServeWebsocket starts a http server and registers a websocket handler
func ListenAndServeWebsocket(addr, path string, handler netutil.ConnHandler, async bool) error {
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
	if async {
		ln, err := net.Listen("tcp", httpServer.Addr)
		if err != nil {
			return err
		}
		go httpServer.Serve(netutil.NewTCPKeepAliveListener(ln.(*net.TCPListener), time.Minute*3))
		return nil
	}
	return httpServer.ListenAndServe()
}
