package service

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/google/uuid"
	"github.com/gopherd/log"

	"github.com/gopherd/doge/build"
	"github.com/gopherd/doge/config"
	"github.com/gopherd/doge/erron"
	"github.com/gopherd/doge/mq"
	"github.com/gopherd/doge/os/signal"
	"github.com/gopherd/doge/service/component"
	"github.com/gopherd/doge/service/discovery"
	"github.com/gopherd/doge/time/timer"
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

// ConfigRewriter rewrites Config
type ConfigRewriter interface {
	RewriteConfig(unsafe.Pointer)
}

var (
	startedAt time.Time
	pid       int
)

func init() {
	startedAt = time.Now()
	pid = os.Getpid()
}

// Since returns duration since started time
func Since() time.Duration {
	return time.Since(startedAt)
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
	log.Info().Print("service started, Ctrl-C or run command `kill -s INT <pid> to shutdown the service`")
	signal.Listen()
	log.Info().Print("service received signal INT")
	app.SetState(Stopping)

	if app.Busy() {
		log.Info().Print("service busy now, maybe waiting for a while")
		ticker := time.NewTicker(time.Millisecond * 100)
		defer ticker.Stop()
		for range ticker.C {
			if !app.Busy() {
				break
			}
		}
	}

	log.Info().Print("shutting down service")
	app.SetState(Closed)

	return app.Shutdown()
}

// BaseService implements Service
type BaseService struct {
	self  Service
	name  string
	id    int64
	uuid  string
	state State
	force bool

	config struct {
		ptr       atomic.Value
		canReload bool
	}

	discovery  discovery.Discovery
	mq         mq.Conn
	components *component.Manager

	tickers struct {
		keepalive *timer.Ticker
		reloadCfg *timer.Ticker
	}
}

// NewBaseService creates a BaseService
func NewBaseService(self Service, cfg config.Configurator) *BaseService {
	s := &BaseService{
		self:       self,
		uuid:       strings.ReplaceAll(uuid.NewString(), "-", ""),
		components: component.NewManager(),
	}
	s.config.ptr.Store(cfg)
	_, s.config.canReload = self.(ConfigRewriter)
	s.tickers.keepalive = timer.NewTicker(time.Second * 3)
	s.tickers.reloadCfg = timer.NewTicker(time.Second * 2)
	return s
}

func (app *BaseService) AddComponent(com component.Component) component.Component {
	return app.components.Add(com)
}

// Name implements Service Name method
func (app *BaseService) Name() string {
	return app.name
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

// MQ returns the mq engine
func (app *BaseService) MQ() mq.Conn {
	return app.mq
}

func (app *BaseService) register(nx bool) error {
	if app.discovery == nil {
		return nil
	}
	var content DiscoveryContent
	now := time.Now().UnixNano() / 1e6
	content.State.Updated = now
	content.State.PID = pid
	content.State.State = app.state
	cfg := app.config.ptr.Load().(config.Configurator)
	if d, ok := cfg.(config.Discoverable); ok {
		content.Config = d.DiscoveredContent()
	} else {
		content.Config = cfg
	}
	data, err := json.Marshal(content)
	if err != nil {
		return err
	}
	id := strconv.FormatInt(app.ID(), 10)
	err = app.discovery.Register(context.Background(), app.Name(), id, string(data), nx)
	if err != nil {
		if discovery.IsExist(err) {
			loaded, err := app.discovery.Find(context.Background(), app.Name(), id)
			if err != nil {
				return err
			}
			var old DiscoveryContent
			if err := json.Unmarshal([]byte(loaded), &old); err != nil {
				return err
			}
			closed := old.State.State == Closed
			expired := old.State.Updated+2*int64(app.tickers.keepalive.Interval()/time.Millisecond) < now
			if !closed {
				if expired {
					if !app.force {
						log.Error().
							String("name", app.Name()).
							Int64("id", app.ID()).
							String("uuid", app.uuid).
							Print("service found and not closed but expired, you can startup with command line flag -F")
						return erron.New("service not closed")
					}
					log.Warn().
						String("name", app.Name()).
						Int64("id", app.ID()).
						Print("force startup service")
				} else {
					log.Warn().
						String("name", app.Name()).
						Int64("id", app.ID()).
						Print("service not closed, stop it first!")
					return erron.New("service not closed")
				}
			}
			if err := app.discovery.Unregister(context.Background(), app.Name(), id); err != nil {
				return err
			}
			return app.discovery.Register(context.Background(), app.Name(), id, string(data), nx)
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
	return app.discovery.Unregister(context.Background(), app.Name(), strconv.FormatInt(app.ID(), 10))
}

func (app *BaseService) reloadCfg() {
	if !app.config.canReload {
		return
	}
	defer func() {
		if e := recover(); e != nil {
			log.Error().Any("error", e).Print("reload config panicked")
		}
	}()

	cfg := app.config.ptr.Load().(config.Configurator)
	rewriter, ok := app.self.(ConfigRewriter)
	if !ok {
		app.config.canReload = false
		return
	}
	newCfg := cfg.Default()
	newCfg.SetSource(cfg.GetSource())
	if err := config.Read(newCfg, true); err != nil {
		log.Warn().Error("error", err).Print("reload config error")
		return
	}
	newCfg.OnReload()
	rewriter.RewriteConfig(unsafe.Pointer(reflect.ValueOf(newCfg).Pointer()))
	app.config.ptr.Store(newCfg)
	log.Trace().Print("reload config successfully")
}

// Init implements Service Init method
func (app *BaseService) Init() error {
	cfg := app.config.ptr.Load().(config.Configurator)
	defaultSource := build.Name() + ".conf"
	flag.CommandLine.BoolVar(&app.force, "F", false, "Whether force startup service while it not closed but expired")
	err := config.Init(flag.CommandLine, cfg, config.WithDefaultSource(defaultSource))
	if err != nil {
		return erron.Throw(err)
	}
	core := cfg.GetCore()
	app.id = core.ID
	app.name = core.Name
	if app.ID() <= 0 {
		return erron.Throwf("invalid service id: %d", app.ID())
	}
	if app.Name() == "" {
		return erron.Throwf("invalid service name: %q", app.Name())
	}

	// initialize log
	level, ok := log.ParseLevel(core.Log.Level)
	if !ok {
		level = log.LevelInfo
	}
	prefix := core.Log.Prefix
	if prefix == "" {
		prefix = build.Name()
	}
	var (
		options []log.Option
		writers []log.Writer
	)
	for _, source := range core.Log.Writers {
		w, err := log.Open(source)
		if err != nil {
			return erron.Throwf("open writer %q error: %s", source, err.Error())
		}
		writers = append(writers, w)
	}
	if len(writers) == 0 {
		options = append(options, log.WithOutput(os.Stderr))
	} else {
		options = append(options, log.WithWriters(writers...))
	}
	options = append(options, log.WithLevel(level))
	options = append(options, log.WithPrefix(prefix))
	options = append(options, log.WithFlags(core.Log.FixedFlags()))
	log.Start(options...)
	log.Info().
		Int("pid", pid).
		String("name", app.Name()).
		Int64("id", app.ID()).
		String("uuid", app.uuid).
		Print("initializing service")

	// open discovery
	if !core.Discovery.Off {
		d, err := discovery.Open(core.Discovery.Name, core.Discovery.Source)
		if err != nil {
			return erron.Throw(err)
		}
		app.discovery = d
	}

	// open mq
	if !core.MQ.Off {
		if app.discovery == nil {
			return erron.Throwf("discovery required if mq enabled")
		}
		q, err := mq.Open(core.MQ.Name, core.MQ.Source, app.discovery)
		if err != nil {
			return erron.Throw(err)
		}
		app.mq = q
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
	if app.tickers.keepalive.Next(now) {
		app.register(false)
	}
	if app.tickers.reloadCfg.Next(now) {
		app.reloadCfg()
	}
}
