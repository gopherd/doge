package config

import (
	"encoding/json"
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
	var s string
	switch mode {
	case Dev:
		s = "dev"
	case Preview:
		s = "preview"
	case Prod:
		s = "prod"
	default:
		return nil, fmt.Errorf("unknown mode: %d", mode)
	}
	return json.Marshal(s)
}

// UnmarshalJSON implements json.Unmarshaler UnmarshalJSON method
func (mode *Mode) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
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
	Core CoreConfig `json:"core"`
}

// Read implements Configurator Read method
func (c *BaseConfig) Read(self Configurator, r io.Reader) error {
	return jsonx.NewDecoder(r,
		jsonx.WithSupportComment(),
		jsonx.WithSupportExtraComma(),
		jsonx.WithSupportUnquotedKey(),
	).Decode(self)
}

// Write implements Configurator Write method
func (c BaseConfig) Write(self Configurator, w io.Writer) error {
	var enc = jsonx.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "\t")
	return enc.Encode(self)
}

// Core configuration
type CoreConfig struct {
	Project   string          `json:"project"`
	Mode      Mode            `json:"mode"`
	ID        int64           `json:"id"`
	Log       LogConfig       `json:"log"`
	MQ        MQConfig        `json:"mq"`
	Discovery DiscoveryConfig `json:"discovery"`
}

// SetSource implements Configurator SetSource method
func (c *BaseConfig) SetSource(source string) {
	c.source = source
}

// GetCore implements Configurator GetCore method
func (c *BaseConfig) GetCore() *CoreConfig {
	return &c.Core
}

// LogConfig represents configuration of log
type LogConfig struct {
	// Level of log
	Level string `json:"level"`
	// Prefix to preappend to each log message
	Prefix string `json:"prefix"`
	// Flags of log printer, @see githug.com/mkideal/log@Flags.
	// -1: no flags
	//  0: default flags
	Flags int `json:"flags"`

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
