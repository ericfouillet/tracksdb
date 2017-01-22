package main

import (
	"github.com/ericfouillet/tracksdb"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	tracksdb.Start()
}
