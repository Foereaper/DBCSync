// Copyright (c) 2025 DBCsync
//
// DBCsync is licensed under the MIT License.
// See the LICENSE file for details.

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
            jobCfg := cfg.SyncJobs["item"]
            return syncItemTemplate(conns, jobCfg)
        },
    },
}