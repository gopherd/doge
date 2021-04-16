package config

import (
	"errors"
	"fmt"
	"io"
	"net/url"

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

var (
	ErrPasswordRequired = errors.New("password required")
	ErrSchemeMismatched = errors.New("scheme mismatched")
)

type URLConfig interface {
	MarshalURL() (*url.URL, error)
	UnmarshalURL(*url.URL) error
}

func MarshalURL(cfg URLConfig) (string, error) {
	u, err := cfg.MarshalURL()
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func UnmarshalURL(cfg URLConfig, rawurl string) error {
	u, err := url.Parse(rawurl)
	if err != nil {
		return err
	}
	return cfg.UnmarshalURL(u)
}

// RedisConfig represents redis client config
type RedisConfig struct {
	Host     string
	Port     int
	User     string
	Password string

	Timeout int
}

func (r RedisConfig) MarshalURL() (*url.URL, error) {
	u := &url.URL{
		Scheme: "redis",
		Host:   fmt.Sprintf("%s:%d", r.Host, r.Port),
	}
	if r.User != "" {
		if r.Password != "" {
			u.User = url.UserPassword(r.User, r.Password)
		} else {
			// password required if user represented
			return nil, ErrPasswordRequired
		}
	} else {
		// assign password to username if user represented
		u.User = url.User(r.Password)
	}

	return u, nil
}

func (r *RedisConfig) UnmarshalURL(u *url.URL) error {
	if u.Scheme != "redis" {
		return ErrSchemeMismatched
	}

	return nil
}
