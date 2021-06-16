package discovery

import (
	"context"
	"fmt"
	"sync"
)

// Discovery represents a interface for service discovery
type Discovery interface {
	// Register registers a service
	Register(ctx context.Context, name, id, content string) error
	// Unregister unregisters a service
	Unregister(ctx context.Context, name, id string) error
	// Resolve resolves any one service by name
	Resolve(ctx context.Context, name string) (id, content string, err error)
	// Resolve resolves all services by name
	ResolveAll(ctx context.Context, name string) (map[string]string, error)
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
func Register(driverName string, driver Driver) {
	if driver == nil {
		panic("discovery: Register driver is nil")
	}
	driversMu.Lock()
	defer driversMu.Unlock()
	if _, dup := drivers[driverName]; dup {
		panic("discovery: Register called twice for driver " + driverName)
	}
	drivers[driverName] = driver
}

// Open opens a discovery specified by its discovery driver name and
// a driver-specific source.
func Open(driverName, source string) (Discovery, error) {
	driversMu.RLock()
	driver, ok := drivers[driverName]
	driversMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("discovery: unknown driver %q (forgotten import?)", driverName)
	}
	return driver.Open(source)
}
