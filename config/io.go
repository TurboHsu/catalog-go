package config

import (
	"log"
	"os"

	"github.com/pelletier/go-toml/v2"
)

func (c *Config) Load(path string) (err error) {
	// Check if file exists
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			*c = DefaultConfig()
			if e := c.Save(path); e != nil {
				log.Fatalf("[F] Failed to save default config: %v\n", e)
			}
			log.Printf("[I] Default config saved to %s, exiting.\n", path)
			os.Exit(0)
		}
		log.Fatalf("[F] Failed to load config: %v\n", err)
	}
			

	buf, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	toml.Unmarshal(buf, c)
	return nil
}

func (c *Config) Save(path string) (err error) {
	buf, err := toml.Marshal(c)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, buf, 0644)
	if err != nil {
		return err
	}
	return nil
}