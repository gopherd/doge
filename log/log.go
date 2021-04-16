package log

const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
)

// global printer
var gprinter = newStdPrinter()

type startOptions struct {
	enableHTTPHandler bool
	printer           Printer
	writers           []Writer
}

func (opt *startOptions) apply(options []Option) {
	for i := range options {
		if options[i] != nil {
			options[i](opt)
		}
	}
}

// Option is option for Start
type Option func(*startOptions)

// WithHTTPHandler enable or disable http handler for settting level
func WithHTTPHandler(yes bool) Option {
	return func(opt *startOptions) {
		opt.enableHTTPHandler = yes
	}
}

// WithPrinter specify custom printer
func WithPrinter(printer Printer) Option {
	return func(opt *startOptions) {
		opt.printer = printer
	}
}

// WithPrinter appends a custom writer
func WithWriter(writer Writer) Option {
	return func(opt *startOptions) {
		opt.writers = append(opt.writers, writer)
	}
}

// WithConsle append a console writer
func WithConsle() Option {
	return WithWriter(newConsole())
}

// WithFile append a file writer
func WithFile(fileOptions FileOptions) Option {
	return WithWriter(newFile(fileOptions))
}

// WithMultiFile append a multifile writer
func WithMultiFile(multiFileOptions MultiFileOptions) Option {
	return WithWriter(newMultiFile(multiFileOptions))
}

// Start inits global printer with options
func Start(options ...Option) error {
	var opt startOptions
	opt.apply(options)
	if opt.printer != nil && len(opt.writers) > 0 {
		println("log.Start: writers ignored because printer sepecfied")
	}
	if opt.printer == nil {
		switch len(opt.writers) {
		case 0:
			return nil
		case 1:
			opt.printer = newPrinter(opt.writers[0], true)
		default:
			opt.printer = newPrinter(mixWriter{opt.writers}, true)
		}
	}
	gprinter.Shutdown()
	gprinter = opt.printer
	gprinter.Start()
	if opt.enableHTTPHandler {
		registerHTTPHandlers()
	}
	return nil
}

// Shutdown shutdowns global printer
func Shutdown() {
	gprinter.Shutdown()
}

// GetLevel gets level of global printer
func GetLevel() Level {
	return gprinter.GetLevel()
}

// SetLevel sets level of gloabl printer
func SetLevel(level Level) {
	gprinter.SetLevel(level)
}

// Trace prints log with trace level
func Trace(format string, args ...interface{}) {
	gprinter.Printf(1, LvTRACE, format, args...)
}

// Debug prints log with debug level
func Debug(format string, args ...interface{}) {
	gprinter.Printf(1, LvDEBUG, format, args...)
}

// Info prints log with info level
func Info(format string, args ...interface{}) {
	gprinter.Printf(1, LvINFO, format, args...)
}

// Warn prints log with warning level
func Warn(format string, args ...interface{}) {
	gprinter.Printf(1, LvWARN, format, args...)
}

// Error prints log with error level
func Error(format string, args ...interface{}) {
	gprinter.Printf(1, LvERROR, format, args...)
}

// Trace prints log with fatal level
func Fatal(format string, args ...interface{}) {
	gprinter.Printf(1, LvFATAL, format, args...)
}

// Printf wraps global printer Printf method
func Printf(calldepth int, level Level, format string, args ...interface{}) {
	gprinter.Printf(calldepth, level, format, args...)
}

// logger implements Logger
type logger struct {
	prefix string
}

func NewLogger(prefix string) Logger {
	return logger{prefix: prefix}
}

func (l logger) prependPrefix(format string) string {
	if len(l.prefix) > 0 {
		return l.prefix + format
	} else {
		return format
	}
}

func (l logger) Trace(format string, args ...interface{}) {
	gprinter.Printf(1, LvTRACE, l.prependPrefix(format), args...)
}

func (l logger) Debug(format string, args ...interface{}) {
	gprinter.Printf(1, LvDEBUG, l.prependPrefix(format), args...)
}

func (l logger) Info(format string, args ...interface{}) {
	gprinter.Printf(1, LvINFO, l.prependPrefix(format), args...)
}

func (l logger) Warn(format string, args ...interface{}) {
	gprinter.Printf(1, LvWARN, l.prependPrefix(format), args...)
}

func (l logger) Error(format string, args ...interface{}) {
	gprinter.Printf(1, LvERROR, l.prependPrefix(format), args...)
}

func (l logger) Fatal(format string, args ...interface{}) {
	gprinter.Printf(1, LvFATAL, l.prependPrefix(format), args...)
}
