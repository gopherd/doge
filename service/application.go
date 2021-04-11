package service

import (
	"flag"
	"os"

	"github.com/gopherd/doge/build"
	"github.com/gopherd/doge/config"
	"github.com/gopherd/doge/osutil/signal"
)

// Application represents a process
type Application interface {
	Service
	ID() int                           // ID of Application
	Configurator() config.Configurator // Config of application
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

	if err := app.Init(); err != nil {
		return err
	}
	app.Start()

	// Waiting signal INT, you can kill the process via
	//
	//	kill -s INT <pid>
	//
	// or Ctrl-C
	signal.Register(os.Interrupt, func(os.Signal) bool {
		return true
	})
	signal.Listen()

	return app.Shutdown()
}

// BaseApplication implments Application
type BaseApplication struct {
	id   int
	name string
}

// NewBaseApplication creates a BaseApplication
func NewBaseApplication() *BaseApplication {
	return &BaseApplication{}
}

// SetID sets id of application
func (app *BaseApplication) SetID(id int) { app.id = id }

// SetName sets name of application
func (app *BaseApplication) SetName(name string) { app.name = name }

// ID implements Application ID method
func (app *BaseApplication) ID() int {
	return app.id
}

// Name implements Application Name method
func (app *BaseApplication) Name() string {
	return app.name
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
