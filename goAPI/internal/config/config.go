package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Port int `yaml:"port"`
}

type HMailConfig struct {
	AdminUser string `yaml:"admin_user"`
	AdminPass string `yaml:"admin_pass"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
	HMail  HMailConfig  `yaml:"hmail"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析 YAML 失败: %v", err)
	}

	return &cfg, nil
}