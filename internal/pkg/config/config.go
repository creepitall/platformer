package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

type Config struct {
	Title       string
	WindowsSize WindowSize
	VSync       bool
	EnableFPS   bool
	FPS         time.Duration
	Logging     bool
}

type WindowSize struct {
	X float64
	Y float64
}

func (c *Config) ReadFile(filepath string) {
	type aliasWindowSize struct {
		X float64 `yaml:"x"`
		Y float64 `yaml:"y"`
	}
	type alias struct {
		Title       string          `yaml:"title"`
		WindowsSize aliasWindowSize `yaml:"windows_size"`
		VSync       bool            `yaml:"vsync"`
		EnableFPS   bool            `yaml:"enable_fps"`
		FPS         int             `yaml:"fps"`
		Logging     bool            `yaml:"logging"`
	}
	var tmp alias

	yamlFile, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Error reading config file: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, &tmp)
	if err != nil {
		fmt.Printf("Error unmarshall yaml file: %v", err)
	}

	c.Title = tmp.Title
	c.VSync = tmp.VSync
	c.WindowsSize = WindowSize{
		X: tmp.WindowsSize.X,
		Y: tmp.WindowsSize.Y}
	c.EnableFPS = tmp.EnableFPS
	c.FPS = time.Duration(tmp.FPS)
	c.Logging = tmp.Logging
}

func LoadConfig() *Config {
	cfg := &Config{}
	cfg.ReadFile("cfg/values.yaml")

	return cfg
}
