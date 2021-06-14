package component

import (
	"reflect"
	"time"

	"github.com/mkideal/log"
)

// Component represents a generic logic component
type Component interface {
	// Name returns name of component
	Name() string
	// Init initializes the component
	Init() error
	// Start starts the component
	Start()
	// Shutdown shutdwons the component
	Shutdown()
	// Update updates the component per frame
	Update(time.Time, time.Duration)
	// Logger returns a logger interface
	Logger() log.Logger
}

// BaseComponent implements the Component interface{}
type BaseComponent struct {
	name   string
	logger log.Logger
}

// NewBaseComponent creates a BaseComponent
func NewBaseComponent(name string) *BaseComponent {
	return &BaseComponent{
		name:   name,
		logger: log.NewLogger(name),
	}
}

// Name implements Component Name method
func (com *BaseComponent) Name() string {
	return com.name
}

// Init implements Component Init method
func (com *BaseComponent) Init() error {
	return nil
}

// Start implements Component Start method
func (com *BaseComponent) Start() {
}

// Shutdown implements Component Shutdown method
func (com *BaseComponent) Shutdown() {
}

// Update implements Component Update method
func (com *BaseComponent) Update(now time.Time, dt time.Duration) {
}

// Logger implements Component Logger method
func (com *BaseComponent) Logger() log.Logger {
	return com.logger
}

// Manager used to manages a group of components
type Manager struct {
	components      []Component
	type2components map[reflect.Type]Component
}

// NewManager creates a Manager
func NewManager() *Manager {
	return &Manager{
		type2components: make(map[reflect.Type]Component),
	}
}

// Add adds a component to the manager
func (m *Manager) Add(com Component) Component {
	t := reflect.TypeOf(com).Elem()
	if _, found := m.type2components[t]; found {
		panic("component type " + t.String() + " duplicated")
	}
	m.components = append(m.components, com)
	m.type2components[t] = com
	return com
}

// Find finds a component from the manager by type
func (m *Manager) Find(t reflect.Type) Component {
	return m.type2components[t]
}

// Len returns the number of components
func (m *Manager) Len() int {
	return len(m.components)
}

// Get returns ith component
func (m *Manager) Get(i int) Component {
	return m.components[i]
}

// Init initializes all components
func (m *Manager) Init() error {
	for i := range m.components {
		m.components[i].Logger().Info("component initializing")
		if err := m.components[i].Init(); err != nil {
			m.components[i].Logger().Info("component initialize error: %v", err)
			return err
		}
		m.components[i].Logger().Info("component initialized")
	}
	return nil
}

// Start starts all components
func (m *Manager) Start() {
	for i := range m.components {
		m.components[i].Logger().Info("component starting")
		m.components[i].Start()
		m.components[i].Logger().Info("component started")
	}
}

// Shutdown shutdowns all components in reverse order
func (m *Manager) Shutdown() {
	for i := len(m.components) - 1; i >= 0; i-- {
		m.components[i].Logger().Info("component shutting down")
		m.components[i].Shutdown()
		m.components[i].Logger().Info("component shutted down")
	}
}

// Update updates all components
func (m *Manager) Update(now time.Time, dt time.Duration) {
	for i := range m.components {
		m.components[i].Update(now, dt)
	}
}
