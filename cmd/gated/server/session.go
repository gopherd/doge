package server

import (
	"net"
	"sync/atomic"
	"time"

	"github.com/gopherd/doge/jwt"
	"github.com/gopherd/doge/net/netutil"
)

type user struct {
	device string
	token  jwt.Payload
}

type state int

const (
	stateCreated state = iota
	statePendingLogin
	stateLogged
	stateClosing
	stateOverflow
)

type handler interface {
	onReady(*session)
	onClose(*session, error)
	onMessage(*session, netutil.Body) error
}

type session struct {
	*netutil.Session
	id                int64
	ip                string
	state             int32
	user              user
	createdAt         int64
	lastKeepaliveTime int64
	lastUpdateSidTime int64
	currSceneId       int64

	handler handler
}

func newSession(id int64, ip string, conn net.Conn, handler handler) *session {
	s := &session{
		id:        id,
		ip:        ip,
		state:     int32(stateCreated),
		createdAt: time.Now().UnixNano() / 1e6,
		handler:   handler,
	}
	s.Session = netutil.NewSession(conn, s)
	return s
}

// OnReady implements netutil.SessionEventHandler OnReady method
func (s *session) OnReady() {
	s.handler.onReady(s)
}

// OnClose implements netutil.SessionEventHandler OnClose method
func (s *session) OnClose(err error) {
	s.handler.onClose(s, err)
}

// OnMessage implements netutil.SessionEventHandler OnMessage method
func (s *session) OnMessage(body netutil.Body) error {
	return s.handler.onMessage(s, body)
}

func (s *session) getState() state {
	return state(atomic.LoadInt32(&s.state))
}

func (s *session) setState(state state) {
	atomic.StoreInt32(&s.state, int32(state))
}

func (s *session) getUser() user {
	return s.user
}

func (s *session) setUser(user user) {
	s.user = user
}
