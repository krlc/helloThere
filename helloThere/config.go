package helloThere

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

type Config struct {
	InfluxDB struct {
		Address     string
		Token       string
		Database    string
		Measurement string
	}

	GeoIPFile string

	Common struct {
		BatchTime    time.Duration
		ResetTime    time.Duration
		MaxUsers     uint
		CookieOnly   bool
		IPFilterOnly bool
		BehindProxy  bool
	}

	Endpoint struct {
		Address string
		Cert    string
		Key     string
		Data    []byte
		MIME    string
	}
}

func NewConfig(path string) (*Config, error) {
	config := new(Config)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}
