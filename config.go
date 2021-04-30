package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pelletier/go-toml"
)

type config struct {
	Mail  mailConfig  `toml:"mail"`
	Drive watchConfig `toml:"drive"`
}

type mailConfig struct {
	// Human-readable name of the email sernder.
	Sender string `toml:"sender"`

	// Basic authentication

	Username string `toml:"username"`
	Password string `toml:"password"`

	// Mail Server

	Server string `toml:"server"`
	Port   uint   `toml:"port"`

	// Who should receive the email.
	RecepientList []string `toml:"recepient_list"`

	// Language of the email (this does not affect loggine).
	Language string `toml:"lang"`
}

type watchConfig struct {
	Path       string `toml:"path"`
	LimitBytes int64  `toml:"limit_bytes"`
}

func configFromFile(configFile string) (*config, error) {
	cfgData, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file %q: %v", configFile, err)
	}
	cfg := config{}
	if err := toml.Unmarshal(cfgData, &cfg); err != nil {
		return nil, fmt.Errorf("cannot unmarshal config file %q: %v", configFile, err)
	}
	return &cfg, nil
}

func (c *config) validate(acceptableLanguages []string) error {
	for _, v := range acceptableLanguages {
		if c.Mail.Language == v {
			return nil
		}
	}
	return fmt.Errorf("lang must be one of: (%s)", strings.Join(acceptableLanguages, ", "))
}
