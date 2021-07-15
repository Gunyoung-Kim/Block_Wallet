package main

import (
	"github.com/Gunyoung-Kim/blockchain/explorer"
	"github.com/Gunyoung-Kim/blockchain/rest"
)

func main() {
	go explorer.Start(3000)
	rest.Start(4000)
}
