package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Work       time.Duration `json:"work"`
	ShortBreak time.Duration `json:"short_break"`
	LongBreak  time.Duration `json:"long_break"`
	Cycle      int           `json:"cycle"`
}

func DefaultConfig() *Config {
	return &Config{
		Work:       25 * time.Minute,
		ShortBreak: 5 * time.Minute,
		LongBreak:  15 * time.Minute,
		Cycle:      4,
	}
}

func Load(filename string) *Config {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			config := DefaultConfig()
			config.Save(filename)
			return config
		}
		return DefaultConfig()
	}
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Printf("Ha ocurrido un error en el fichero: %v", err)
		return DefaultConfig()
	}

	return &config
}

func (c *Config) Save(filename string) error {
	if err := c.Validate(); err != nil {
		return err
	}
	jsonData, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return fmt.Errorf("error al hacer marshal: %w", err)
	}
	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("ha ocurrido un error: %w", err)
	}

	return nil
}

func (c *Config) Validate() error {
	if c.Work <= 0 || c.LongBreak <= 0 || c.ShortBreak <= 0 {
		return fmt.Errorf("las duraciones deben ser positivas")
	}
	if c.Cycle <= 0 {
		return fmt.Errorf("el ciclo debe ser positivo")
	}
	return nil
}
