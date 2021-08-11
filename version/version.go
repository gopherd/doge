package version

import (
	_ "embed"
	"expvar"
	"strings"
)

//go:embed VERSION
var version string

func NoDuplicate(name string) {
	expvar.Publish(name, expvar.NewString(strings.TrimSpace(version)))
}
