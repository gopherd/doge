package server

import (
	"sync"

	"github.com/gopherd/doge/component"
)

type pendingSession struct {
	uid  int64
	meta uint32
}

type sessionManager struct {
	*component.BaseComponent

	maxConns      int
	maxConnsPerIP int

	context interface {
		GetConfig() *Config
	}

	mutex    sync.RWMutex
	sessions map[int64]*session
	uid2sid  map[int64]int64
	ips      map[string]int
}

func newSessionManager(server *server) *sessionManager {
	return &sessionManager{
		BaseComponent: component.NewBaseComponent("session_manager"),
		sessions:      make(map[int64]*session),
		uid2sid:       make(map[int64]int64),
		ips:           make(map[string]int),
	}
}

// Init overrides BaseComponent Init method
func (c *sessionManager) Init() error {
	if err := c.BaseComponent.Init(); err != nil {
		return err
	}
	cfg := c.context.GetConfig()
	c.maxConns = cfg.MaxConns
	c.maxConnsPerIP = cfg.MaxConnsPerIP
	return nil
}

func (c *sessionManager) size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return len(c.sessions)
}

func (c *sessionManager) add(s *session) (n int, ok bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.sessions[s.id] = s
	n = len(c.sessions)
	if n < c.maxConns {
		ok = true
	} else {
		s.setState(stateOverflow)
	}
	return
}

func (c *sessionManager) remove(id int64) *session {
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
	if uid := s.getUid(); uid > 0 {
		delete(c.uid2sid, uid)
	}
	return s
}
