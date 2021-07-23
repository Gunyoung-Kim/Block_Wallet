package main

import (
	"github.com/Gunyoung-Kim/blockchain/cli"
	"github.com/Gunyoung-Kim/blockchain/db"
)

func main() {
	defer db.Close()
	cli.Start()
}
