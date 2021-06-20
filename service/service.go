package service

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/mkideal/log"

	"github.com/gopherd/doge/build"
	"github.com/gopherd/doge/config"
	"github.com/gopherd/doge/erron"
	"github.com/gopherd/doge/os/signal"
	"github.com/gopherd/doge/service/component"
	"github.com/gopherd/doge/service/discovery"
)

// State represents service state
type State int

const (
	Closed   State = iota // Closed service
	Running               // Running service
	Stopping              // Stopping service
)

func (state State) String() string {
	switch state {
	case Closed:
		return "Closed"
	case Running:
		return "Running"
	case Stopping:
		return "Stopping"
	default:
		return "Unknown(" + strconv.Itoa(int(state)) + ")"
	}
}

// DiscoveryContent is the content for registering
type DiscoveryContent struct {
	Config interface{} `json:"config"` // config of service
	State  struct {
		PID   int   `json:"pid"`   // process id
		State State `json:"state"` // run state
	} `json:"state"` // runtime state of service
}

// Service represents a process
type Service interface {
	// ID returns id of service
	ID() int64
	// Name of service
	Name() string
	// SetState sets state of service
	SetState(state State) error
	// Busy reports whether the service is busy
	Busy() bool
	// Init initializes the service
	Init() error
	// Start starts the service
	Start() error
	// Shutdown shutdowns the service
	Shutdown() error
}

// Run runs the service
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

var pid = os.Getpid()

func exec(app Service) error {
	log.Start(
		log.WithConsle(),
		log.WithLevel(log.LvDEBUG),
		log.WithPrefix(build.Name()),
	)
	defer log.Shutdown()

	log.Info().Int("pid", pid).Print("initializing service")
	if err := app.Init(); err != nil {
		return err
	}
	log.Info().Print("starting service")
	app.Start()
	if err := app.SetState(Running); err != nil {
		log.Info().Error("error", err).Print("set service state error")
		app.Shutdown()
		return err
	}

	// Waiting signal INT, you can kill the process via
	//
	//	kill -s INT <pid>
	//
	// or Ctrl-C
	signal.Register(os.Interrupt, func(os.Signal) bool {
		return true
	})
	log.Info().Print("service started, waiting signal INT")
	signal.Listen()
	log.Info().Print("service received signal INT")
	app.SetState(Stopping)

	if app.Busy() {
		log.Info().Print("service now busy, waiting...")
		ticker := time.NewTicker(time.Millisecond * 100)
		defer ticker.Stop()
		for range ticker.C {
			if !app.Busy() {
				break
			}
		}
	}

	app.SetState(Closed)
	log.Info().Print("shutting down service")
	return app.Shutdown()
}

// BaseService implements Service
type BaseService struct {
	name       string
	id         int64
	state      State
	cfg        config.Configurator
	discovery  discovery.Discovery
	components *component.Manager
}

// NewBaseService creates a BaseService
func NewBaseService(cfg config.Configurator) *BaseService {
	return &BaseService{
		cfg:        cfg,
		components: component.NewManager(),
	}
}

func (app *BaseService) AddComponent(com component.Component) component.Component {
	return app.components.Add(com)
}

// SetName sets name of service
func (app *BaseService) SetName(name string) {
	app.name = name
}

// Name implements Service Name method
func (app *BaseService) Name() string {
	return app.name
}

// SetID sets id of service
func (app *BaseService) SetID(id int64) {
	app.id = id
}

// ID implements Service ID method
func (app *BaseService) ID() int64 {
	return app.id
}

// SetState implements Service SetState method
func (app *BaseService) SetState(state State) error {
	app.state = state
	return app.register()
}

// State returns state of service
func (app *BaseService) State() State {
	return app.state
}

// Busy implements Service Busy method
func (app *BaseService) Busy() bool {
	return false
}

// Discovery returns the discovery engine
func (app *BaseService) Discovery() discovery.Discovery {
	return app.discovery
}

func (app *BaseService) register() error {
	if app.discovery == nil {
		return nil
	}
	var content DiscoveryContent
	content.State.PID = pid
	content.State.State = app.state
	if d, ok := app.cfg.(config.Discoverable); ok {
		content.Config = d.DiscoveredContent()
	} else {
		content.Config = app.cfg
	}
	data, err := json.Marshal(content)
	if err != nil {
		return err
	}
	return app.discovery.Register(
		context.Background(),
		app.name,
		strconv.FormatInt(app.id, 10),
		string(data),
		false,
	)
}

func (app *BaseService) unregister() error {
	if app.discovery == nil {
		return nil
	}
	return app.discovery.Unregister(context.Background(), app.name, strconv.FormatInt(app.id, 10))
}

// Init implements Service Init method
func (app *BaseService) Init() error {
	defaultSource := build.Name() + ".conf"
	err := config.Init(flag.CommandLine, app.cfg, config.WithDefaultSource(defaultSource))
	if err != nil {
		return erron.Throw(err)
	}
	name, source := app.cfg.GetDiscovery()
	if name != "" {
		d, err := discovery.Open(name, source)
		if err != nil {
			return erron.Throw(err)
		}
		app.discovery = d
	}
	return app.components.Init()
}

// Start implements Service Start method
func (app *BaseService) Start() {
	app.components.Start()
}

// Shutdown implements Service Shutdown method
func (app *BaseService) Shutdown() {
	app.components.Shutdown()
	app.unregister()
}

// Update updates per frame
func (app *BaseService) Update(now time.Time, dt time.Duration) {
	app.components.Update(now, dt)
}
