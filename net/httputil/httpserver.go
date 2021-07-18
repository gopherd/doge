package httputil

import (
	"context"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gopherd/log"
	xnetutil "golang.org/x/net/netutil"

	"github.com/gopherd/doge/net/netutil"
)

var pong = []byte{'p', 'o', 'n', 'g'}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(pong)
}

type Config struct {
	Address           string `json:"address"`
	StaticDir         string `json:"static_dir"`
	StaticPath        string `json:"static_path"`
	ConnTimeout       int64  `json:"conn_timeout"`
	ReadHeaderTimeout int64  `json:"read_header_timeout"`
	ReadTimeout       int64  `json:"read_timeout"`
	WriteTimeout      int64  `json:"write_timeout"`
	MaxConns          int    `json:"max_conns"`
}

func (cfg *Config) autofix() {
	if cfg.ConnTimeout <= 0 {
		cfg.ConnTimeout = 60 // 1 分钟
	}
	if cfg.ReadHeaderTimeout <= 0 {
		cfg.ReadHeaderTimeout = 30
	}
	if cfg.ReadTimeout <= 0 {
		cfg.ReadTimeout = 30
	}
	if cfg.WriteTimeout <= 0 {
		cfg.WriteTimeout = 30
	}
	if cfg.MaxConns <= 0 {
		cfg.MaxConns = 4096
	}
}

// HTTPServer ...
type HTTPServer struct {
	cfg    Config
	mux    http.ServeMux
	server *http.Server

	numHandling int64
}

func NewHTTPServer(cfg Config) *HTTPServer {
	cfg.autofix()
	httpd := &HTTPServer{
		cfg: cfg,
	}
	httpd.server = &http.Server{
		Addr:              httpd.cfg.Address,
		Handler:           &httpd.mux,
		ReadHeaderTimeout: time.Duration(httpd.cfg.ReadHeaderTimeout) * time.Second,
		ReadTimeout:       time.Duration(httpd.cfg.ReadTimeout) * time.Second,
		WriteTimeout:      time.Duration(httpd.cfg.WriteTimeout) * time.Second,
	}
	return httpd
}

func (httpd *HTTPServer) NumHandling() int64 {
	return atomic.LoadInt64(&httpd.numHandling)
}

func (httpd *HTTPServer) ListenAndServe(async bool) error {
	addr := httpd.server.Addr
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	l := netutil.NewTCPKeepAliveListener(ln.(*net.TCPListener), time.Duration(httpd.cfg.ConnTimeout)*time.Second)
	l2 := xnetutil.LimitListener(l, httpd.cfg.MaxConns)
	if async {
		go func() {
			if err := httpd.server.Serve(l2); err != nil {
				if err == http.ErrServerClosed {
					log.Info().Print("http server closed")
				} else {
					log.Error().Error("error", err).Print("http server")
				}
			}
		}()
		return nil
	} else {
		return httpd.server.Serve(l2)
	}
}

func (httpd *HTTPServer) Shutdown(ctx context.Context) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
	}
	return httpd.server.Shutdown(ctx)
}

func (httpd *HTTPServer) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request), middlewares ...Middleware) {
	httpd.Handle(pattern, http.HandlerFunc(handler), middlewares...)
}

func (httpd *HTTPServer) Handle(pattern string, handler http.Handler, middlewares ...Middleware) {
	for _, m := range middlewares {
		handler = m.Apply(handler)
	}
	httpd.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&httpd.numHandling, 1)
		defer atomic.AddInt64(&httpd.numHandling, -1)
		w.Header().Add("Connection", "Keep-alive")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Keep-alive", "30")
		handler.ServeHTTP(w, r)
	})
}

func (httpd *HTTPServer) JSONResponse(w http.ResponseWriter, r *http.Request, data interface{}, options ...ResponseOptions) {
	JSONResponse(w, data, options...)
}
