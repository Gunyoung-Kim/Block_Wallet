package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
)

// Block is component of block chain
type Block struct {
	Height   int    `json:"height"`
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
}

type blockChain struct {
	blocks []*Block
}

var b *blockChain
var once sync.Once

func (b *Block) calculateHash() {
	hash := sha256.Sum256([]byte(b.Data + b.PrevHash))
	b.Hash = fmt.Sprintf("%x", hash)
}

func getLastHash() string {
	totalBlocks := len(GetBlockChain().blocks)
	if totalBlocks == 0 {
		return ""
	}

	return GetBlockChain().blocks[totalBlocks-1].Hash
}

func createBlock(data string) *Block {
	newBlock := Block{Data: data, Hash: "", PrevHash: getLastHash(), Height: len(GetBlockChain().blocks) + 1}
	newBlock.calculateHash()
	return &newBlock
}

func (b *blockChain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
}

func (b *blockChain) AllBlocks() []*Block {
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

var ErrNotFound = errors.New("Block Not Found")

func (b *blockChain) GetBlock(height int) (*Block, error) {

	if len(b.blocks) < height || height-1 < 0 {
		return nil, ErrNotFound
	}

	return b.blocks[height-1], nil
}
