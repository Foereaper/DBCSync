module dbcsync

go 1.22.0

require (
    github.com/go-sql-driver/mysql v1.9.3
)

replace github.com/go-sql-driver/mysql => ../dep/mysql
