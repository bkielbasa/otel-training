package tracescleaner

import (
	"errors"
	"regexp"
)

type Config struct {
	Exclude []string `mapstructure:"exclude"`
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
