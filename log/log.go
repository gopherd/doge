package log

const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
)

// global printer
var gprinter = newStdPrinter()

// Shutdown shutdowns global printer
func Shutdown() {
	gprinter.Shutdown()
}

// InitWithPrinter inits global printer with a specified printer
func Start(p Printer) error {
	gprinter.Shutdown()
	gprinter = p
	gprinter.Start()
	return nil
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
