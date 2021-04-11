package service

// Service represents a runnable unit
type Service interface {
	// Name returns name of service
	Name() string
	// Init initializes the service
	Init() error
	// Start starts the service
	Start() error
	// Shutdown shutdown the service
	Shutdown() error
}

// Manager manages a group services
type Manager struct {
	services []Service
}

// NewManager creates a manager
func NewManager() *Manager {
	return &Manager{}
}

// Add adds a service to manager
func (m *Manager) Add(service Service) {
	m.services = append(m.services, service)
}

// Init initializes all services
func (m *Manager) Init() error {
	for i := range m.services {
		if err := m.services[i].Init(); err != nil {
			return err
		}
	}
	return nil
}

// Start starts all services
func (m *Manager) Start() error {
	for i := range m.services {
		if err := m.services[i].Start(); err != nil {
			return err
		}
	}
	return nil
}

// Shutdown shutdowns all services
func (m *Manager) Shutdown() error {
	for i := range m.services {
		if err := m.services[i].Shutdown(); err != nil {
			return err
		}
	}
	return nil
}
