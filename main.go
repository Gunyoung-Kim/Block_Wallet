package main

import (
	"fmt"

	"github.com/Gunyoung-Kim/blockchain/blockchain"
)

func main() {
	chain := blockchain.GetBlockChain()
	chain.AddBlock("Second Block")
	chain.AddBlock("Thrid Block")
	chain.AddBlock("Fourth Block")
	for _, block := range chain.AllBlocks() {
		fmt.Printf("data : %s\n", block.Data())
		fmt.Printf("hash: %s\n", block.Hash())
		fmt.Printf("prevHash: %s\n", block.PrevHash())
	}
}
