package module

import (
	"reflect"
	"time"

	"github.com/gopherd/log"
)

// Module represents a generic logic module
type Module interface {
	// Name returns the name of module
	Name() string
	// Init initializes the module
	Init() error
	// Start starts the module
	Start()
	// Shutdown shutdwons the module
	Shutdown()
	// Update updates the module per frame
	Update(time.Time, time.Duration)
}

// BaseModule implements the Module interface{}
type BaseModule struct {
	name   string
	logger log.Prefix
}

// NewBaseModule creates a BaseModule
func NewBaseModule(name string) *BaseModule {
	return &BaseModule{
		name:   name,
		logger: log.Prefix(name),
	}
}

// Name implements Module Name method
func (com *BaseModule) Name() string {
	return com.name
}

// Init implements Module Init method
func (com *BaseModule) Init() error {
	return nil
}

// Start implements Module Start method
func (com *BaseModule) Start() {
}

// Shutdown implements Module Shutdown method
func (com *BaseModule) Shutdown() {
}

// Update implements Module Update method
func (com *BaseModule) Update(now time.Time, dt time.Duration) {
}

// Logger returns the Module logger
func (com *BaseModule) Logger() log.Prefix {
	return com.logger
}

// Manager used to manages a group of modules
type Manager struct {
	modules      []Module
	type2modules map[reflect.Type][]Module
}

// NewManager creates a Manager
func NewManager() *Manager {
	return &Manager{
		type2modules: make(map[reflect.Type][]Module),
	}
}

// Add adds a module to the manager
func (m *Manager) Add(com Module) Module {
	t := reflect.TypeOf(com).Elem()
	m.type2modules[t] = append(m.type2modules[t], com)
	m.modules = append(m.modules, com)
	return com
}

// Find finds the first added module from the manager by type
func (m *Manager) Find(t reflect.Type) Module {
	coms, ok := m.type2modules[t]
	if !ok || len(coms) == 0 {
		return nil
	}
	return coms[0]
}

// FindAll finds all modules from the manager by type
func (m *Manager) FindAll(t reflect.Type) []Module {
	return m.type2modules[t]
}

// Len returns the number of modules
func (m *Manager) Len() int {
	return len(m.modules)
}

// Get returns ith module
func (m *Manager) Get(i int) Module {
	return m.modules[i]
}

// Init initializes all modules
func (m *Manager) Init() error {
	for _, com := range m.modules {
		name := com.Name()
		log.Prefix(name).Info().Print("module initializing")
		if err := com.Init(); err != nil {
			log.Prefix(name).Info().Error("error", err).Print("module initialize error")
			return err
		}
		log.Prefix(name).Info().Print("module initialized")
	}
	return nil
}

// Start starts all modules
func (m *Manager) Start() {
	for _, com := range m.modules {
		name := com.Name()
		log.Prefix(name).Info().Print("module starting")
		com.Start()
		log.Prefix(name).Info().Print("module started")
	}
}

// Shutdown shutdowns all modules in reverse order
func (m *Manager) Shutdown() {
	for i := len(m.modules) - 1; i >= 0; i-- {
		com := m.modules[i]
		name := com.Name()
		log.Prefix(name).Info().Print("module shutting down")
		com.Shutdown()
		log.Prefix(name).Info().Print("module shutted down")
	}
}

// Update updates all modules
func (m *Manager) Update(now time.Time, dt time.Duration) {
	for i := range m.modules {
		m.modules[i].Update(now, dt)
	}
}
