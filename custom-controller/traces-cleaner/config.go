package tracescleaner

import (
	"errors"
	"regexp"

	"go.opentelemetry.io/collector/receiver/otlpreceiver"
)

type Config struct {
	Exclude   []string               `mapstructure:"exclude"`
	OTEL      otlpreceiver.Config    `mapstructure:"otel"`
	Protocols otlpreceiver.Protocols `mapstructure:"protocols"`
}

func (cfg *Config) Validate() error {
	if len(cfg.Exclude) == 0 {
		return errors.New("exclude list is empty")
	}

	for _, e := range cfg.Exclude {
		_, err := regexp.Compile(e)
		if err != nil {
			return err
		}
	}

	return nil
}
