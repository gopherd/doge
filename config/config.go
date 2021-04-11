package config

import (
	"io"

	"github.com/gopherd/doge/encoding/jsonx"
)

// Config implments Configurator
type Config struct {
	// source of config
	source string `json:"_"`
}

// Read implements Configurator Read method
func (c *Config) Read(self Configurator, r io.Reader) error {
	var dec = jsonx.NewDecoder(r)
	return dec.Decode(self)
}

// Write implements Configurator Write method
func (c Config) Write(self Configurator, w io.Writer) error {
	var enc = jsonx.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "\t")
	return enc.Encode(self)
}

// SetSource implements Configurator SetSource method
func (c *Config) SetSource(source string) {
	c.source = source
}
