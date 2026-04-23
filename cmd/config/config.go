package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type (
	// Config holds environment configurations.
	Config struct {
		DB       DB
		Migrator Migrator
		Service  Service
	}

	// DB holds the environment variables for a DB setup.
	DB struct {
		Host     string `env:"DB_HOST"`
		Port     string `env:"DB_PORT"`
		Username string `env:"DB_USERNAME"`
		Password string `env:"DB_PASSWORD"`
		Database string `env:"DB_DATABASE"`
	}

	// Migrator holds the environment variables for a DB migrator.
	Migrator struct {
		Path string `env:"DB_MIGRATOR_PATH"`
	}

	// Service holds the environment variables for a service setup.
	Service struct {
		Port string `env:"SERVICE_PORT"`
	}
)

// Setup initializes config and loads environment variables.
func Setup() (*Config, error) {
	var cfg = &Config{}

	if err := loadEnv(cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func loadEnv(cfg any) error {
	rv := reflect.ValueOf(cfg)

	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("unabled to load env, config is not a none-nil pointer")
	}

	v := rv.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		if !field.CanSet() {
			continue
		}

		if field.Kind() == reflect.Struct {
			if err := loadEnv(field.Addr().Interface()); err != nil {
				return err
			}

			continue
		}

		var tag string

		if tag = t.Field(i).Tag.Get("env"); tag == "" {
			continue
		}

		val, exists := os.LookupEnv(tag)
		if !exists {
			return errors.New(fmt.Sprintf("missing environment variable: %s", tag))
		}

		if field.CanSet() {
			field.SetString(strings.TrimSpace(val))
		}
	}

	return nil
}
