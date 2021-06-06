package config

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/gopherd/doge/build"
)

// Configurator represents generic config of application
type Configurator interface {
	// Read reads config from reader r
	Read(Configurator, io.Reader) error
	// Write writes config to writer w
	Write(Configurator, io.Writer) error
	// SetSource sets source of config
	SetSource(string)
	// GetDiscovery returns configuration of discovery
	GetDiscovery() (name, source string)
}

type option struct {
	defaultSource             string
	errorHandling             func(error)
	flagInput, usageInput     string
	flagOutput, usageOutput   string
	flagVersion, usageVersion string
}

func newOption() *option {
	return &option{
		flagInput:    "c",
		flagOutput:   "e",
		usageOutput:  "Exported config filename",
		flagVersion:  "v",
		usageVersion: "Print version information",
	}
}

func (opt option) getUsageInput() string {
	if opt.usageInput != "" {
		return opt.usageInput
	}
	if opt.defaultSource != "" {
		return "Config source (default " + opt.defaultSource + ")"
	}
	return "Config source"
}

func (opt option) getUsageOutput() string {
	return opt.usageOutput
}

func (opt option) getUsageVersion() string {
	return opt.usageVersion
}

func (opt *option) applyOptions(options []Option) {
	for _, o := range options {
		o(opt)
	}
}

// Option is the option for Init
type Option func(*option)

// WithDefaultSource specify default config source
func WithDefaultSource(source string) Option {
	return func(opt *option) {
		opt.defaultSource = source
	}
}

// WithInput specify command line flag name and usage for config input
func WithInput(flag, usage string) Option {
	return func(opt *option) {
		opt.flagInput = flag
		opt.usageInput = usage
	}
}

// WithOutput specify command line flag name and usage for config output
func WithOutput(flag, usage string) Option {
	return func(opt *option) {
		opt.flagOutput = flag
		opt.usageOutput = usage
	}
}

// WithVersion specify command line flag name and usage for version
func WithVersion(flag, usage string) Option {
	return func(opt *option) {
		opt.flagVersion = flag
		opt.usageVersion = usage
	}
}

type exitError struct {
	code int
}

func (e exitError) Error() string {
	return fmt.Sprintf("exit with code %d", e.code)
}

func IsExitError(err error) (code int, ok bool) {
	if err == nil {
		return 0, false
	}
	if e, ok := err.(exitError); ok {
		return e.code, ok
	}
	return 0, false
}

// Init initializes Configure cfg from command line flags with options
func Init(flagSet *flag.FlagSet, cfg Configurator, options ...Option) error {
	var opt = newOption()
	opt.applyOptions(options)

	var (
		input, output      string
		shouldPrintVersion bool
	)
	flagSet.StringVar(&input, opt.flagInput, "", opt.getUsageInput())
	flagSet.StringVar(&output, opt.flagOutput, "", opt.getUsageOutput())
	flagSet.BoolVar(&shouldPrintVersion, opt.flagVersion, false, opt.getUsageVersion())
	flagSet.Parse(os.Args[1:])

	if shouldPrintVersion {
		build.Print()
		return exitError{0}
	}

	var optional = false
	if input == "" && opt.defaultSource != "" {
		optional = true
		input = opt.defaultSource
	}

	if inputFile, err := os.Open(input); err != nil {
		if !optional || !os.IsNotExist(err) {
			return err
		}
	} else {
		cfg.SetSource(input)
		err := cfg.Read(cfg, inputFile)
		inputFile.Close()
		if err != nil {
			return err
		}
	}

	if output != "" {
		if outputFile, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666); err != nil {
			return err
		} else {
			err := cfg.Write(cfg, outputFile)
			outputFile.Close()
			if err != nil {
				return err
			}
		}
		return exitError{0}
	}

	return nil
}

// Discoverable represents a discoverable configurator
type Discoverable interface {
	// DiscoveredContent returns a discovered data
	DiscoveredContent() ([]byte, error)
}
