package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type block struct {
	data     string
	hash     string
	prevHash string
}

type blockChain struct {
	blocks []*block
}

var b *blockChain
var once sync.Once

func (b block) Data() string {
	return b.data
}

func (b block) Hash() string {
	return b.hash
}

func (b block) PrevHash() string {
	return b.prevHash
}

func (b *block) calculateHash() {
	hash := sha256.Sum256([]byte(b.data + b.prevHash))
	b.hash = fmt.Sprintf("%x", hash)
}

func getLastHash() string {
	totalBlocks := len(GetBlockChain().blocks)
	if totalBlocks == 0 {
		return ""
	}

	return GetBlockChain().blocks[totalBlocks-1].hash
}

func createBlock(data string) *block {
	newBlock := block{data: data, hash: "", prevHash: getLastHash()}
	newBlock.calculateHash()
	return &newBlock
}

func (b *blockChain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
}

func (b *blockChain) AllBlocks() []*block {
	return b.blocks
}

// GetBlockChain get blockChain
func GetBlockChain() *blockChain {
	if b == nil {
		once.Do(func() {
			b = &blockChain{}
			b.AddBlock("Genesis Block")
		})
	}
	return b
}
