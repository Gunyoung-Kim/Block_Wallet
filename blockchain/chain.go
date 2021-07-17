package blockchain

import (
	"sync"

	"github.com/Gunyoung-Kim/blockchain/db"
	"github.com/Gunyoung-Kim/blockchain/utils"
)

type blockChain struct {
	NewestHash string `json:"newestHash"`
	Height     int    `json:"height"`
}

var b *blockChain // variable for singleton pattern of blockChain
var once sync.Once

//------------ receiver function for blockChain ------------------

func (b *blockChain) persist() {
	db.SaveCheckPoint(utils.ToBytes(b))
}

func (b *blockChain) restoreFromBytes(data []byte) {
	utils.FromBytes(b, data)
}

//AddBlock createBlock using current NewestHash and Height
// and update NewestHash and Height for blockChain
func (b *blockChain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.persist()
}

//Blocks return all pointer of Blocks from DB
func (b *blockChain) Blocks() []*Block {
	hashCursor := b.NewestHash
	var result []*Block
	for {
		block, _ := FindBlock(hashCursor)
		result = append(result, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return result
}

//------------ function for blockChain ------------------

// BlockChain get blockChain
// This function is for singleton pattern of blockChain
func BlockChain() *blockChain {
	if b == nil {
		once.Do(func() {
			b = &blockChain{"", 0}
			checkPoint := db.CheckPoint()
			if checkPoint == nil {
				b.AddBlock("Genesis Block")
			} else {
				b.restoreFromBytes(checkPoint)
			}
		})
	}
	return b
}
