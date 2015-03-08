package config

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
)

type Config struct {
	data    map[string]interface{}
	changed bool
}

func (c *Config) Set(key string, value interface{}) {
	orig, found := c.data[key]
	if found {
		if !reflect.DeepEqual(value, orig) {
			c.changed = true
		}
	} else {
		c.changed = true
	}
	c.data[key] = value
}

func (c *Config) GetString(key string, defaultValue string) string {
	v, found := c.data[key]
	if !found {
		return defaultValue
	}

	if vs, ok := v.(string); ok {
		return vs
	}
	return defaultValue
}

func (c *Config) GetInt(key string, defaultValue int) int {
	v, found := c.data[key]
	if !found {
		return defaultValue
	}

	if vi, ok := v.(float64); ok {
		return int(vi)
	}
	return defaultValue
}

func New() *Config {
	c := new(Config)
	c.data = make(map[string]interface{})
	return c
}

var ErrEmptyFile = errors.New("empty file")

func From(name string) (*Config, error) {
	c := New()
	f, err := os.OpenFile(name, os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// check content size
	st, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if st.Size() == 0 {
		return nil, ErrEmptyFile
	}

	if err = json.NewDecoder(f).Decode(&c.data); err != nil {
		return nil, err
	}
	return c, nil
}

func SaveTo(name string, c *Config) error {
	if !c.changed {
		return nil
	}

	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(&c.data)
}
