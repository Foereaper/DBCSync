module dbcsync

go 1.22.0

require github.com/go-sql-driver/mysql v1.9.3

require filippo.io/edwards25519 v1.1.0 // indirect

replace github.com/go-sql-driver/mysql => ../dep/mysql
