// Copyright (c) 2025 DBCsync
//
// DBCsync is licensed under the MIT License.
// See the LICENSE file for details.

package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// DBConfig holds config for a single database
type DBConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Name     string `json:"name"`
}

type SyncJobConfig struct {
	MinEntry int `json:"minEntry"`
}

type Config struct {
	World    DBConfig                 `json:"world"`
	DBC      DBConfig                 `json:"dbc"`
	SyncJobs map[string]SyncJobConfig `json:"syncJobs"`
}

// loadOrInitConfig loads config.json, or generates a template if missing
func loadOrInitConfig(path string) (*Config, bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Create template config
		template := Config{
            World: DBConfig{"root", "password", "127.0.0.1", "3306", "world"},
            DBC:   DBConfig{"root", "password", "127.0.0.1", "3306", "dbc"},
            SyncJobs: map[string]SyncJobConfig{
                "item": {MinEntry: 1000},
            },
        }

		data, err := json.MarshalIndent(template, "", "  ")
		if err != nil {
			return nil, false, fmt.Errorf("marshal template: %w", err)
		}

		if err := os.WriteFile(path, data, 0644); err != nil {
			return nil, false, fmt.Errorf("write template: %w", err)
		}

		return nil, true, nil
	}

	// Load existing config
	file, err := os.Open(path)
	if err != nil {
		return nil, false, fmt.Errorf("open config: %w", err)
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, false, fmt.Errorf("decode config: %w", err)
	}
	return &cfg, false, nil
}