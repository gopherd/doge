package discovery

import (
	"fmt"
	"sync"
)

// Discovery represents a interface for service discovery
type Discovery interface {
	// Register registers a service
	Register(serviceName, serviceId string, value []byte) error
	// Unregister unregisters a service
	Unregister(serviceName, serviceId string) error
	// Resolve resolves any one service by name
	Resolve(serviceName string) ([]byte, error)
	// Resolve resolves all services by name
	ResolveAll(serviceName string) (map[string][]byte, error)
}

// Driver is the interface that must be implemented by a discovery driver
type Driver interface {
	// Open returns a new discovery instance by a driver-specific source name
	Open(source string) (Discovery, error)
}

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]Driver)
)

// Register makes a discovery driver available by the provided name
func Register(name string, driver Driver) {
	if driver == nil {
		panic("discovery: Register driver is nil")
	}
	driversMu.Lock()
	defer driversMu.Unlock()
	if _, dup := drivers[name]; dup {
		panic("discovery: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

// Open opens a discovery specified by its discovery driver name and
// a driver-specific source name.
func Open(name string, source string) (Discovery, error) {
	driversMu.RLock()
	driver, ok := drivers[name]
	driversMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("discovery: unknown driver %q (forgotten import?)", name)
	}
	return driver.Open(source)
}
