package event

// Type represents type of event
type Type int64

// Event represents an event data
type Event interface {
	Type() Type
}

// Handler handles event
type Handler interface {
	Handle(Event) (swallowed bool)
}

// HandlerFunc wraps a function as Handler
type HandlerFunc func(Event) (swallowed bool)

// Handle implements Handler HandleEvent method
func (fn HandlerFunc) Handle(e Event) bool { return fn(e) }

// Registry manages event handlers
type Registry struct {
	handlers map[Type][]Handler
}

// NewRegistry creates an event registry
func NewRegistry() *Registry {
	return &Registry{
		handlers: make(map[Type][]Handler),
	}
}

// Handle registers event handler by type
func (r *Registry) Handle(t Type, h Handler) {
	r.handlers[t] = append(r.handlers[t], h)
}

// HandleFunc registers event handler func by type
func (r *Registry) HandleFunc(t Type, h HandlerFunc) {
	r.handlers[t] = append(r.handlers[t], h)
}

// Post posts event to handlers by type
func (r *Registry) Post(e Event) {
	handlers, ok := r.handlers[e.Type()]
	if !ok {
		return
	}
	for i := range handlers {
		if handlers[i].Handle(e) {
			return
		}
	}
}
