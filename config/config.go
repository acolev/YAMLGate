package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Gateway  GatewayConfig `yaml:"gateway"`
	Services []Service     `yaml:"services"`
	Settings Settings      `yaml:"settings"`
}

type GatewayConfig struct {
	Address       string   `yaml:"address"`
	GlobalHeaders []Header `yaml:"headers"`
}

type Service struct {
	Name          string   `yaml:"name"`
	ProxyURL      string   `yaml:"proxy_url"`
	Headers       []Header `yaml:"headers"`
	Routes        []Route  `yaml:"routes"`
	NeedsChromedp bool     `yaml:"chromedp"`
}

type Route struct {
	ServicePath string `yaml:"service_path"`
	GatewayPath string `yaml:"gateway_path"`
	Method      string `yaml:"method"`
}

type Header struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type Settings struct {
	Timeout              string `yaml:"timeout"`
	MaxRetries           int    `yaml:"max_retries"`
	CacheDefaultDuration string `yaml:"cache_default_duration"`
}

// LoadConfig загружает конфигурацию из YAML файла
func LoadConfig(configPath string) (*Config, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать конфигурационный файл: %w", err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("не удалось распарсить YAML: %w", err)
	}

	return &cfg, nil
}
