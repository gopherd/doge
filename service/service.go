package service

import (
	"encoding/json"
	"flag"
	"os"
	"time"

	"github.com/mkideal/log"

	"github.com/gopherd/doge/build"
	"github.com/gopherd/doge/config"
	"github.com/gopherd/doge/os/signal"
	"github.com/gopherd/doge/service/component"
	"github.com/gopherd/doge/service/discovery"
)

// Service represents a process
type Service interface {
	// ID returns id of service
	ID() string
	// Name of service
	Name() string
	// Init initializes the service
	Init() error
	// Start starts the service
	Start() error
	// Shutdown shutdowns the service
	Shutdown() error
}

// Run runs the application
func Run(app Service) {
	if err := exec(app); err != nil {
		code, ok := config.IsExitError(err)
		if !ok {
			code = 1
		}
		if code != 0 {
			println(err.Error())
		}
		os.Exit(code)
	}
}

func exec(app Service) error {
	log.Start(
		log.WithConsle(),
		log.WithLevel(log.LvDEBUG),
		log.WithPrefix(build.Name()),
	)
	defer log.Shutdown()

	log.Info("initializing application, pid = %d", os.Getpid())
	if err := app.Init(); err != nil {
		return err
	}
	log.Info("starting application")
	app.Start()

	// Waiting signal INT, you can kill the process via
	//
	//	kill -s INT <pid>
	//
	// or Ctrl-C
	signal.Register(os.Interrupt, func(os.Signal) bool {
		return true
	})
	log.Info("application started, waiting signal INT")
	signal.Listen()
	log.Info("application received signal INT")

	log.Info("shutting down application")
	return app.Shutdown()
}

// BaseService implements Service
type BaseService struct {
	id         string
	name       string
	cfg        config.Configurator
	discovery  discovery.Discovery
	components *component.Manager
}

// NewBaseService creates a BaseApplication
func NewBaseService(cfg config.Configurator) *BaseService {
	return &BaseService{
		cfg:        cfg,
		components: component.NewManager(),
	}
}

func (app *BaseService) AddComponent(com component.Component) component.Component {
	return app.components.Add(com)
}

// SetID sets id of application
func (app *BaseService) SetID(id string) {
	app.id = id
}

// ID implements Application.Service ID method
func (app *BaseService) ID() string {
	return app.id
}

// SetName sets name of application
func (app *BaseService) SetName(name string) {
	app.name = name
}

// Name implements Application.Service Name method
func (app *BaseService) Name() string {
	return app.name
}

// Init implements Application.Service Init method
func (app *BaseService) Init() error {
	defaultSource := build.Name() + ".conf"
	err := config.Init(flag.CommandLine, app.cfg, config.WithDefaultSource(defaultSource))
	if err != nil {
		return err
	}
	name, source := app.cfg.GetDiscovery()
	if name != "" {
		d, err := discovery.Open(name, source)
		if err != nil {
			return err
		}
		app.discovery = d
		var discoveredContent []byte
		if dcfg, ok := app.cfg.(config.Discoverable); ok {
			discoveredContent, err = dcfg.DiscoveredContent()
		} else {
			discoveredContent, err = json.Marshal(app.cfg)
		}
		if err != nil {
			return err
		}
		if err := app.discovery.Register(app.name, app.id, discoveredContent); err != nil {
			return nil
		}
	}
	return nil
}

// Start implements Application.Service Start method
func (app *BaseService) Start() {
	app.components.Start()
}

// Shutdown implements Application.Service Shutdown method
func (app *BaseService) Shutdown() {
	app.components.Shutdown()
}

// Update updates per frame
func (app *BaseService) Update(now time.Time, dt time.Duration) {
	app.components.Update(now, dt)
}
