package blockchain

import (
	"sync"

	"github.com/Gunyoung-Kim/blockchain/db"
	"github.com/Gunyoung-Kim/blockchain/utils"
)

const (
	defaultDifficulty   int = 2
	difficultyInterval  int = 5
	blockCreateInterval int = 2
	allowedRange        int = 2
)

type blockChain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDifficulty"`
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
	b.CurrentDifficulty = block.Difficulty
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

func (b *blockChain) recalculateDifficulty() int {
	allBlocks := b.Blocks()
	newestBlock := allBlocks[0]
	lastCheckedBlock := allBlocks[difficultyInterval-1]
	actualTime := (newestBlock.Timestamp / 60) - (lastCheckedBlock.Timestamp / 60)
	expectedTime := difficultyInterval * blockCreateInterval
	if actualTime < (expectedTime - allowedRange) {
		return b.CurrentDifficulty + 1
	} else if actualTime > (expectedTime + allowedRange) {
		return b.CurrentDifficulty - 1
	}
	return b.CurrentDifficulty
}

func (b *blockChain) difficulty() int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		return b.recalculateDifficulty()
	} else {
		return b.CurrentDifficulty
	}
}

//------------ function for blockChain ------------------

// BlockChain get blockChain
// This function is for singleton pattern of blockChain
func BlockChain() *blockChain {
	if b == nil {
		once.Do(func() {
			b = &blockChain{
				Height: 0,
			}
			checkPoint := db.CheckPoint()
			if checkPoint == nil {
				b.AddBlock("Genesis")
			} else {
				b.restoreFromBytes(checkPoint)
			}
		})
	}
	return b
}
