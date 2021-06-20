package service

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
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
		Updated int64  `json:"updated"`
		PID     int    `json:"pid"`   // process id
		State   State  `json:"state"` // run state
		UUID    string `json:"uuid"`  // instance uuid
	} `json:"state"` // runtime state of service
}

// Service represents a process
type Service interface {
	// ID returns id of service
	ID() int64
	// Name of service
	Name() string
	// Global unique instance id of service
	UUID() string
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
	defer log.Shutdown()

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
	name              string
	id                int64
	uuid              string
	state             State
	cfg               config.Configurator
	discovery         discovery.Discovery
	force             bool
	components        *component.Manager
	lastKeepaliveTime time.Time
}

// NewBaseService creates a BaseService
func NewBaseService(cfg config.Configurator) *BaseService {
	return &BaseService{
		uuid:       strings.ReplaceAll(uuid.NewString(), "-", ""),
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

// UUID implements service UUID method
func (app *BaseService) UUID() string {
	return app.uuid
}

// SetState implements Service SetState method
func (app *BaseService) SetState(state State) error {
	app.state = state
	return app.register(state == Running)
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

const keepalive = 5 * 1000

func (app *BaseService) register(nx bool) error {
	if app.discovery == nil {
		return nil
	}
	var content DiscoveryContent
	now := time.Now().UnixNano() / 1e6
	content.State.Updated = now
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
	id := strconv.FormatInt(app.id, 10)
	err = app.discovery.Register(context.Background(), app.name, id, string(data), nx)
	if err != nil {
		if discovery.IsExist(err) {
			loaded, err := app.discovery.Find(context.Background(), app.name, id)
			if err != nil {
				return err
			}
			var old DiscoveryContent
			if err := json.Unmarshal([]byte(loaded), &old); err != nil {
				return err
			}
			closed := old.State.State == Closed
			expired := old.State.Updated+2*keepalive < now
			if !closed {
				if expired {
					if !app.force {
						log.Error().
							String("name", app.name).
							Int64("id", app.id).
							String("uuid", app.uuid).
							Print("service found and not closed but expired, you can startup with command line flag -F")
						return erron.New("service not closed")
					}
					log.Warn().
						String("name", app.name).
						Int64("id", app.id).
						Print("force startup service")
				} else {
					log.Warn().
						String("name", app.name).
						Int64("id", app.id).
						Print("service not closed, stop it first!")
					return erron.New("service not closed")
				}
			}
			if err := app.discovery.Unregister(context.Background(), app.name, id); err != nil {
				return err
			}
			return app.discovery.Register(context.Background(), app.name, id, string(data), nx)
		} else {
			return err
		}
	}
	return nil
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
	flag.CommandLine.BoolVar(&app.force, "F", false, "Whether force startup service while it not closed but expired")
	err := config.Init(flag.CommandLine, app.cfg, config.WithDefaultSource(defaultSource))
	if err != nil {
		return erron.Throw(err)
	}
	if x, ok := app.cfg.(interface{ GetID() int64 }); ok {
		app.id = x.GetID()
	}
	if x, ok := app.cfg.(interface{ GetLog() config.LogConfig }); ok {
		logcfg := x.GetLog()
		level, ok := log.ParseLevel(logcfg.Level)
		if !ok {
			level = log.LvINFO
		}
		prefix := logcfg.Prefix
		if prefix == "" {
			prefix = build.Name()
		}
		var options []log.Option
		// TODO: writer from logcfg.Writers
		options = append(options, log.WithConsle())
		options = append(options, log.WithLevel(level))
		options = append(options, log.WithPrefix(prefix))
		if logcfg.Flags < 0 {
			options = append(options, log.WithFlags(0))
		} else if logcfg.Flags > 0 {
			options = append(options, log.WithFlags(logcfg.Flags))
		}
		log.Start(options...)
	}
	log.Info().
		Int("pid", pid).
		Int64("id", app.id).
		String("uuid", app.uuid).
		Print("initializing service")
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
	if now.Sub(app.lastKeepaliveTime) > time.Duration(keepalive)*time.Millisecond {
		app.lastKeepaliveTime = now
		app.register(false)
	}
}
