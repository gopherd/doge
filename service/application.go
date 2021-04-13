package service

import (
	"flag"
	"os"

	"github.com/gopherd/doge/build"
	"github.com/gopherd/doge/config"
	"github.com/gopherd/doge/log"
	"github.com/gopherd/doge/osutil/signal"
)

// Application represents a process
type Application interface {
	Service
	ID() int                           // ID of application
	Configurator() config.Configurator // Config of application
	Logger() log.Logger                // Logger of application
	SetLogger(log.Logger)
}

// Run runs the application
func Run(app Application) {
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

func exec(app Application) error {
	// Initializing config of application
	if cfg := app.Configurator(); cfg != nil {
		defaultSource := build.Name() + ".json"
		err := config.Init(flag.CommandLine, cfg, config.WithDefaultSource(defaultSource))
		if err != nil {
			return err
		}
	}

	log.Start(log.NewPrinter(log.NewConsole(), true))
	defer log.Shutdown()
	log.SetLevel(log.LvDEBUG)

	logger := app.Logger()
	if logger == nil {
		logger = log.NewLogger("(" + build.Name() + ") ")
		app.SetLogger(logger)
	}

	logger.Info("initializing service, pid = %d", os.Getpid())
	if err := app.Init(); err != nil {
		return err
	}
	logger.Info("starting service")
	app.Start()

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
	return app.Shutdown()
}

// BaseApplication implments Application
type BaseApplication struct {
	id           int
	name         string
	configurator config.Configurator
	logger       log.Logger
}

// NewBaseApplication creates a BaseApplication
func NewBaseApplication() *BaseApplication {
	return &BaseApplication{}
}

// SetID sets id of application
func (app *BaseApplication) SetID(id int) {
	app.id = id
}

// SetName sets name of application
func (app *BaseApplication) SetName(name string) {
	app.name = name
}

// SetConfigurator sets configurator of application
func (app *BaseApplication) SetConfigurator(configurator config.Configurator) {
	app.configurator = configurator
}

// ID implements Application ID method
func (app *BaseApplication) ID() int {
	return app.id
}

// Name implements Application Name method
func (app *BaseApplication) Name() string {
	return app.name
}

// Configurator implements Application Configurator method
func (app *BaseApplication) Configurator() config.Configurator {
	return app.configurator
}

// Logger ...
func (app *BaseApplication) Logger() log.Logger {
	return app.logger
}

func (app *BaseApplication) SetLogger(logger log.Logger) {
	app.logger = logger
}

// Init implements Application Init method
func (app *BaseApplication) Init() error {
	return nil
}

// Start implements Application Start method
func (app *BaseApplication) Start() error {
	return nil
}

// Shutdown implements Application Shutdown method
func (app *BaseApplication) Shutdown() error {
	return nil
}
