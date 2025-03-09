package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Unit      string  `yaml:"unit"`
	Lang      string  `yaml:"lang"`
	Longitude float64 `yaml:"longitude"`
	Latitude  float64 `yaml:"latitude"`
	OWMAPIKey string  `yaml:"owm_api_key"`
}

func (cfg *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw struct {
		LongitudeStr string `yaml:"longitude"`
		LatitudeStr  string `yaml:"latitude"`
		Unit         string `yaml:"unit"`
		Lang         string `yaml:"lang"`
		OWMAPIKey    string `yaml:"owm_api_key"`
	}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	longitude, err := strconv.ParseFloat(raw.LongitudeStr, 64)
	if err != nil || longitude < -180 || longitude > 180 {
		return fmt.Errorf("invalid longitude: %v (must be between -180 and 180)", raw.LatitudeStr)
	}

	latitude, err := strconv.ParseFloat(raw.LatitudeStr, 64)
	if err != nil || latitude < -90 || latitude > 90 {
		return fmt.Errorf("invalid latitude: %v (must be between -90 and 90)", raw.LatitudeStr)
	}
	cfg.Unit = raw.Unit
	cfg.Lang = raw.Lang
	cfg.Longitude = longitude
	cfg.Latitude = latitude
	cfg.OWMAPIKey = raw.OWMAPIKey
	return nil
}

func New(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	//Override OWM API Key with env var
	if ownKey, ok := os.LookupEnv("OWM_API_KEY"); ok {
		cfg.OWMAPIKey = ownKey
	}

	return &cfg, nil
}
