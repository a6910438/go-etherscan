module etherscan-go

go 1.13

require (
	github.com/a6910438/go-logger v0.0.0-20191031041009-3779e17d0503
	github.com/asdine/storm/v3 v3.1.0
	github.com/boltdb/bolt v1.3.1
	github.com/ethereum/go-ethereum v1.9.9
	github.com/inconshreveable/log15 v0.0.0-20180818164646-67afb5ed74ec
	github.com/jarcoal/httpmock v1.0.4 // indirect
	github.com/jinzhu/configor v1.1.1
	github.com/jinzhu/gorm v1.9.11
	github.com/judwhite/go-svc v1.1.2
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mattn/go-isatty v0.0.11 // indirect
	github.com/mattn/go-sqlite3 v1.11.0
	github.com/onrik/ethrpc v1.0.0
	github.com/pkg/errors v0.8.1
	github.com/shopspring/decimal v0.0.0-20191130220710-360f2bc03045
	github.com/tidwall/gjson v1.3.5 // indirect
)

replace golang.org/x/sys => github.com/golang/sys v0.0.0-20191026070338-33540a1f6037
