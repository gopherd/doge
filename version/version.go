package version

import (
	_ "embed"
	"expvar"
	"strings"
)

//go:embed VERSION
var version string

type str string

func (s str) String() string { return string(s) }

func NoDuplicate(name string) {
	expvar.Publish(name, str(strings.TrimSpace(version)))
}
