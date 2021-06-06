package server

import (
	"sync"
	"time"

	"github.com/gopherd/doge/service"
)

const (
	kMaxDurationForPendingSession = 5000
	kHandlePendingSessionInterval = 200
	kCleanDeadSessionInterval     = 60000
	kUserInfoTTLRatio             = 750
)

type pendingSession struct {
	uid  int64
	meta uint32
}

type connections struct {
	mutex        sync.RWMutex
	maxSize      int
	maxSizePerIP int
	sessions     map[int64]*session
	uid2sid      map[int64]int64
	ips          map[string]int
}

func newConnections(maxSize, maxSizePerIP int) *connections {
	return &connections{
		maxSize:      maxSize,
		maxSizePerIP: maxSizePerIP,
		sessions:     make(map[int64]*session),
		uid2sid:      make(map[int64]int64),
		ips:          make(map[string]int),
	}
}

func (c *connections) size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return len(c.sessions)
}

func (c *connections) add(s *session) (n int, ok bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.sessions[s.id] = s
	n = len(c.sessions)
	if n < c.maxSize {
		ok = true
	} else {
		s.setState(stateOverflow)
	}
	return
}

func (c *connections) remove(id int64) *session {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	s, ok := c.sessions[id]
	if !ok {
		return nil
	}
	ip := s.ip
	if n, ok := c.ips[ip]; n > 1 {
		c.ips[ip] = n - 1
	} else if ok {
		delete(c.ips, ip)
	}
	if uid := s.getUser().token.Uid; uid > 0 {
		delete(c.uid2sid, uid)
	}
	return s
}

type server struct {
	*service.BaseService
	config *Config

	quit, wait chan struct{}

	// components list all components of ram
	components struct {
		// more...
	}
}

// New creates gated service
func New() service.Service {
	cfg := NewConfig()
	s := &server{
		BaseService: service.NewBaseService(cfg),
		config:      cfg,
		quit:        make(chan struct{}),
		wait:        make(chan struct{}),
	}

	return s
}

// Init overrides BaseService Init method
func (s *server) Init() error {
	return s.BaseService.Init()
}

// Start overrides BaseService Start method
func (s *server) Start() error {
	s.BaseService.Start()
	go s.run()
	return nil
}

// Shutdown overrides BaseService Shutdown method
func (s *server) Shutdown() error {
	close(s.quit)
	<-s.wait
	s.BaseService.Shutdown()
	return nil
}

// run runs service's main loop
func (s *server) run() {
	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()

	lastUpdatedAt := time.Now()
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			s.onUpdate(now, now.Sub(lastUpdatedAt))
			lastUpdatedAt = now
		case <-s.quit:
			close(s.wait)
			return
		}
	}
}

func (s *server) onUpdate(now time.Time, dt time.Duration) {
	s.BaseService.Update(now, dt)
}
