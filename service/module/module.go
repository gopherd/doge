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
	logger *log.ContextLogger
}

// NewBaseModule creates a BaseModule
func NewBaseModule(name string) *BaseModule {
	return &BaseModule{
		name:   name,
		logger: log.Prefix(nil, name),
	}
}

// Name implements Module Name method
func (mod *BaseModule) Name() string {
	return mod.name
}

// Init implements Module Init method
func (mod *BaseModule) Init() error {
	return nil
}

// Start implements Module Start method
func (mod *BaseModule) Start() {
}

// Shutdown implements Module Shutdown method
func (mod *BaseModule) Shutdown() {
}

// Update implements Module Update method
func (mod *BaseModule) Update(now time.Time, dt time.Duration) {
}

// Logger returns the Module logger
func (mod *BaseModule) Logger() *log.ContextLogger {
	return mod.logger
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
func (m *Manager) Add(mod Module) Module {
	t := reflect.TypeOf(mod).Elem()
	m.type2modules[t] = append(m.type2modules[t], mod)
	m.modules = append(m.modules, mod)
	return mod
}

// Find finds the first added module from the manager by type
func (m *Manager) Find(t reflect.Type) Module {
	mods, ok := m.type2modules[t]
	if !ok || len(mods) == 0 {
		return nil
	}
	return mods[0]
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
	for _, mod := range m.modules {
		name := mod.Name()
		log.Info().String("module", name).Print("module initializing")
		if err := mod.Init(); err != nil {
			log.Info().String("module", name).Error("error", err).Print("module initialize error")
			return err
		}
		log.Info().String("module", name).Print("module initialized")
	}
	return nil
}

// Start starts all modules
func (m *Manager) Start() {
	for _, mod := range m.modules {
		name := mod.Name()
		log.Info().String("module", name).Print("module starting")
		mod.Start()
		log.Info().String("module", name).Print("module started")
	}
}

// Shutdown shutdowns all modules in reverse order
func (m *Manager) Shutdown() {
	for i := len(m.modules) - 1; i >= 0; i-- {
		mod := m.modules[i]
		name := mod.Name()
		log.Info().String("module", name).Print("module shutting down")
		mod.Shutdown()
		log.Info().String("module", name).Print("module shutted down")
	}
}

// Update updates all modules
func (m *Manager) Update(now time.Time, dt time.Duration) {
	for i := range m.modules {
		m.modules[i].Update(now, dt)
	}
}
