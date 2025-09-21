# DBCSync

dbcsync is a Go-based utility for synchronizing tables from a TrinityCore **world database** to a **DBC database**. The tool is modular, supports multiple sync jobs, and is fully configurable via `config.json`.

---

## Features

- Connects to **world** and **DBC** MySQL databases.  
- Modular sync functions: each table sync can have its own function.  
- Supports **job-specific configuration**, including filtering by minimum entry values.  
- Automatically generates a **template `config.json`** on first run.  
- Transactional and safe `REPLACE INTO` operations to update existing DBC entries.  
- Tested with **TrinityCore 3.3.5 world database**.  

---

## Requirements

- Go 1.20 or higher  
- MySQL database for both **world** and **DBC**  
- DBC database populated via [Stoneharry's Spell Editor DBC Import Tool](https://github.com/stoneharry/WoW-Spell-Editor)  
- Network access to both MySQL servers  

---

## Installation

1. Clone the repository:

```bash
git clone https://github.com/Foereaper/dbcsync --recurse
cd dbcsync
```

2. Build the binary:

```bash
cd src
go build -o ../bin/dbcsync.exe
```

---

## Usage

1. **First run** — generates `config.json` template:

```bash
./bin/dbcsync.exe
```

You should see:

```
Config template created at .config.json. Please edit it and re-run.
```

2. Edit `config.json` to set:

- Database credentials for `world` and `dbc`.

Example:

```json
{
  "world": {
    "user": "root",
    "password": "password",
    "host": "127.0.0.1",
    "port": "3306",
    "name": "world"
  },
  "dbc": {
    "user": "root",
    "password": "password",
    "host": "127.0.0.1",
    "port": "3306",
    "name": "dbc"
  },
  "syncJobs": {
    "item": {
      "minEntry": 1000
    },
  }
}
```

3. **Run the sync**:

```bash
./bin/dbcsync.exe
```

- Each registered sync job will be executed in sequence.  
- Progress and synced counts will be printed to the console.  

---

## Project Structure

```
dbcsync/
├─ dep/                  # External dependencies (e.g., go-sql-driver/mysql)
├─ src/
│  ├─ main.go            # Entry point
│  ├─ config.go          # Config loading and template generation
│  ├─ db.go              # Database connection helpers
│  ├─ registry.go        # Sync registry for modular jobs
│  ├─ sync_item.go       # Example: item_template → item sync
│  └─ ...                # Additional sync files
├─ .gitignore
└─ README.md
```

---

## Adding New Sync Jobs

1. Create a new `sync_<table>.go` file.  
2. Implement a function with signature:

```go
func sync<TableName>(conns *DBConnections, cfg SyncJobConfig) error
```

3. Add the function to `syncRegistry` in `registry.go`:

```go
var syncRegistry = []SyncFunc{
    {
        Name: "item_template → item",
        Func: func(conns *DBConnections, cfg *Config) error {
            return syncItemTemplate(conns, cfg.SyncJobs["item"])
        },
    },
}
```

4. Add a default entry in `config.json` template if desired.

---

## Notes

- The DBC database **must exist** and be populated by Stoneharry's tool before running this sync.  
- The tool uses `REPLACE INTO` to safely update existing DBC entries.  
- Tested with **TrinityCore 3.3.5 world database**; other versions may require mapping adjustments.  

---

© 2025 DBCsync is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.