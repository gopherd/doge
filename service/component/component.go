package component

import (
	"reflect"
	"time"

	"github.com/mkideal/log"
)

// Component represents a generic logic component
type Component interface {
	// Name returns the name of component
	Name() string
	// Init initializes the component
	Init() error
	// Start starts the component
	Start()
	// Shutdown shutdwons the component
	Shutdown()
	// Update updates the component per frame
	Update(time.Time, time.Duration)
}

// BaseComponent implements the Component interface{}
type BaseComponent struct {
	name   string
	logger log.Prefix
}

// NewBaseComponent creates a BaseComponent
func NewBaseComponent(name string) *BaseComponent {
	return &BaseComponent{
		name:   name,
		logger: log.Prefix(name),
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

// Logger returns the component logger
func (com *BaseComponent) Logger() log.Prefix {
	return com.logger
}

// Manager used to manages a group of components
type Manager struct {
	components      []Component
	type2components map[reflect.Type][]Component
}

// NewManager creates a Manager
func NewManager() *Manager {
	return &Manager{
		type2components: make(map[reflect.Type][]Component),
	}
}

// Add adds a component to the manager
func (m *Manager) Add(com Component) Component {
	t := reflect.TypeOf(com).Elem()
	m.type2components[t] = append(m.type2components[t], com)
	m.components = append(m.components, com)
	return com
}

// Find finds the first added component from the manager by type
func (m *Manager) Find(t reflect.Type) Component {
	coms, ok := m.type2components[t]
	if !ok || len(coms) == 0 {
		return nil
	}
	return coms[0]
}

// FindAll finds all components from the manager by type
func (m *Manager) FindAll(t reflect.Type) []Component {
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
	for _, com := range m.components {
		name := com.Name()
		log.Prefix(name).Info().Print("component initializing")
		if err := com.Init(); err != nil {
			log.Prefix(name).Info().Error("error", err).Print("component initialize error")
			return err
		}
		log.Prefix(name).Info().Print("component initialized")
	}
	return nil
}

// Start starts all components
func (m *Manager) Start() {
	for _, com := range m.components {
		name := com.Name()
		log.Prefix(name).Info().Print("component starting")
		com.Start()
		log.Prefix(name).Info().Print("component started")
	}
}

// Shutdown shutdowns all components in reverse order
func (m *Manager) Shutdown() {
	for i := len(m.components) - 1; i >= 0; i-- {
		com := m.components[i]
		name := com.Name()
		log.Prefix(name).Info().Print("component shutting down")
		com.Shutdown()
		log.Prefix(name).Info().Print("component shutted down")
	}
}

// Update updates all components
func (m *Manager) Update(now time.Time, dt time.Duration) {
	for i := range m.components {
		m.components[i].Update(now, dt)
	}
}
