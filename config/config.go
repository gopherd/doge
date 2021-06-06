package config

import (
	"fmt"
	"io"
	"strings"

	"github.com/gopherd/doge/encoding/jsonx"
)

// Process running mode
type Mode int

const (
	Dev Mode = iota
	Preview
	Prod
)

// MarshalJSON implements json.Marshaler MarshalJSON method
func (mode Mode) MarshalJSON() ([]byte, error) {
	switch mode {
	case Dev:
		return []byte("dev"), nil
	case Preview:
		return []byte("preview"), nil
	case Prod:
		return []byte("prod"), nil
	default:
		return nil, fmt.Errorf("unknown mode: %d", mode)
	}
}

// UnmarshalJSON implements json.Unmarshaler UnmarshalJSON method
func (mode *Mode) UnmarshalJSON(data []byte) error {
	switch strings.ToLower(string(data)) {
	case "dev":
		*mode = Dev
	case "preview":
		*mode = Preview
	case "prod":
		*mode = Prod
	default:
		return fmt.Errorf("unknown mode: %q", string(data))
	}
	return nil
}

// BaseConfig implments Configurator
type BaseConfig struct {
	// source of config
	source string `json:"-"`

	// Core represents core common fields
	Core struct {
		Mode      Mode            `json:"mode"`
		ID        int             `json:"id"`
		FPS       int             `json:"fps"`
		Log       LogConfig       `json:"log"`
		MQ        MQConfig        `json:"mq"`
		Discovery DiscoveryConfig `json:"discovery"`
	} `json:"core"`
}

// Read implements Configurator Read method
func (c *BaseConfig) Read(self Configurator, r io.Reader) error {
	return jsonx.NewDecoder(r).Decode(self)
}

// Write implements Configurator Write method
func (c BaseConfig) Write(self Configurator, w io.Writer) error {
	var enc = jsonx.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "\t")
	return enc.Encode(self)
}

// SetSource implements Configurator SetSource method
func (c *BaseConfig) SetSource(source string) {
	c.source = source
}

// GetDiscovery implements Configurator GetDiscovery method
func (c *BaseConfig) GetDiscovery() (name, source string) {
	return c.Core.Discovery.Name, c.Core.Discovery.Source
}

// LogConfig represents configuration of log
type LogConfig struct {
	// Level of log
	Level string `json:"level"`
	// Prefix to preappend to each log message
	Prefix string `json:"prefix"`

	// Writers specified multi-writers, like:
	//	[
	//		"console",
	//		"file://path/to/filename?suffix=.txt"
	//	]
	Writers []string `json:"writers"`
}

// MQConfig ...
type MQConfig struct {
	Name   string `json:"name"`
	Source string `json:"source"`
}

// DiscoveryConfig ...
type DiscoveryConfig struct {
	Name   string `json:"name"`
	Source string `json:"source"`
}
