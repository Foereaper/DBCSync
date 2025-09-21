package main

import (
	"fmt"
	"log"
)

func main() {
	const configPath = "./config.json"

	// Load or init config
	cfg, created, err := loadOrInitConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// First run - generate config template and exit
	if created {
		fmt.Printf("Config template created at %s. Please edit it and re-run.\n", configPath)
		return
	}

	// Open DBs
	worldDB, err := openDB(cfg.World)
	if err != nil {
		log.Fatalf("Failed to connect to world DB: %v", err)
	}
	defer worldDB.Close()

	dbcDB, err := openDB(cfg.DBC)
	if err != nil {
		log.Fatalf("Failed to connect to dbc DB: %v", err)
	}
	defer dbcDB.Close()

	fmt.Println("Connected to both World and DBC databases!")

	conns := &DBConnections{
		World: worldDB,
		DBC:   dbcDB,
	}

	// Run all registered syncs
    for _, sync := range syncRegistry {
        fmt.Printf("Running sync: %s\n", sync.Name)
        if err := sync.Func(conns, cfg); err != nil {
            log.Fatalf("Sync %s failed: %v", sync.Name, err)
        }
        fmt.Printf("Completed sync: %s\n", sync.Name)
    }

	fmt.Println("All syncs completed successfully!")
}