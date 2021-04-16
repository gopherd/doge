package service

import (
	"flag"
	"os"

	"github.com/gopherd/doge/build"
	"github.com/gopherd/doge/config"
	"github.com/gopherd/doge/log"
	"github.com/gopherd/doge/osutil/signal"
)

// Service represents a process
type Service interface {
	ID() int      // ID of service
	Name() string // Name returns name of service

	// Init initializes the service
	Init() error
	// Start starts the service
	Start() error
	// Shutdown shutdown the service
	Shutdown() error
}

// Run runs the service
func Run(service Service) {
	if err := exec(service); err != nil {
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

type configurable interface {
	Configurator() config.Configurator // Config of service
}

type loggerGetter interface {
	Logger() log.Logger
}

type loggerSetter interface {
	SetLogger(log.Logger)
}

func exec(service Service) error {
	// Initializing config of service if it's a configurableService
	if cs, ok := service.(configurable); ok {
		if cfg := cs.Configurator(); cfg != nil {
			defaultSource := build.Name() + ".json"
			err := config.Init(flag.CommandLine, cfg, config.WithDefaultSource(defaultSource))
			if err != nil {
				return err
			}
		}
	}

	log.Start(log.WithConsle())
	defer log.Shutdown()
	log.SetLevel(log.LvDEBUG)

	// Get logger of service or create a new logger
	var logger log.Logger
	if s, ok := service.(loggerGetter); ok {
		logger = s.Logger()
	}
	if logger == nil {
		logger = log.NewLogger("(" + build.Name() + ") ")
		// Try set logger if service implements loggerSetter
		if s, ok := service.(loggerSetter); ok {
			s.SetLogger(logger)
		}
	}

	logger.Info("initializing service, pid = %d", os.Getpid())
	if err := service.Init(); err != nil {
		return err
	}
	logger.Info("starting service")
	service.Start()

	// Waiting signal INT, you can kill the process via
	//
	//	kill -s INT <pid>
	//
	// or Ctrl-C
	signal.Register(os.Interrupt, func(os.Signal) bool {
		return true
	})
	logger.Info("service started, waiting signal INT")
	signal.Listen()
	logger.Info("service received signal INT")

	logger.Info("shutting down service")
	return service.Shutdown()
}

// BaseService implements Service, configurable, loggerGetter, loggerSetter
type BaseService struct {
	id           int
	name         string
	configurator config.Configurator
	logger       log.Logger
}

// NewBaseService creates a BaseService
func NewBaseService() *BaseService {
	return &BaseService{}
}

// SetID sets id of service
func (service *BaseService) SetID(id int) {
	service.id = id
}

// SetName sets name of service
func (service *BaseService) SetName(name string) {
	service.name = name
}

// SetConfigurator sets configurator of service
func (service *BaseService) SetConfigurator(configurator config.Configurator) {
	service.configurator = configurator
}

// ID implements Service ID method
func (service *BaseService) ID() int {
	return service.id
}

// Name implements Service Name method
func (service *BaseService) Name() string {
	return service.name
}

// Configurator implements Service Configurator method
func (service *BaseService) Configurator() config.Configurator {
	return service.configurator
}

// Logger gets logger of service
func (service *BaseService) Logger() log.Logger {
	return service.logger
}

// SetLogger sets logger of service
func (service *BaseService) SetLogger(logger log.Logger) {
	service.logger = logger
}

// Init implements Service Init method
func (service *BaseService) Init() error {
	return nil
}

// Start implements Service Start method
func (service *BaseService) Start() error {
	return nil
}

// Shutdown implements Service Shutdown method
func (service *BaseService) Shutdown() error {
	return nil
}
