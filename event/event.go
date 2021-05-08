package event

// Type represents type of event
type Type int64

// Event represents an event data
type Event interface {
	EventType() Type
}

// Handler handles event
type Handler interface {
	HandleEvent(Event)
}

// HandlerFunc wraps a function as Handler
type HandlerFunc func(Event)

// HandleEvent implements Handler HandleEvent method
func (fn HandlerFunc) HandleEvent(e Event) { fn(e) }

// Manager manages event handlers
type Manager struct {
	handlers map[Type][]Handler
}

// NewManager creates an event manager
func NewManager() *Manager {
	return &Manager{
		handlers: make(map[Type][]Handler),
	}
}

// Register registers handler by event type
func (m *Manager) Register(t Type, h Handler) {
	if handlers, ok := m.handlers[t]; ok {
		m.handlers[t] = append(handlers, h)
	} else {
		m.handlers[t] = []Handler{h}
	}
}

// Post posts event
func (m *Manager) Post(e Event) {
	handlers, ok := m.handlers[e.EventType()]
	if !ok {
		return
	}
	for i := range handlers {
		handlers[i].HandleEvent(e)
	}
}
