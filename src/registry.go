package main

// SyncFunc represents a sync function with a name and handler
type SyncFunc struct {
	Name string
	Func func(conns *DBConnections, cfg *Config) error
}

// syncRegistry holds all sync functions
var syncRegistry = []SyncFunc{
	{
		Name: "item_template → item",
		Func: func(conns *DBConnections, cfg *Config) error {
			return syncItemTemplate(conns, &cfg.SyncJobs["item"])
		},
	},
	// Add more syncs here
}