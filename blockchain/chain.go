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

func (b *blockChain) restoreFromBytes(data []byte) {
	utils.FromBytes(b, data)
}

//AddBlock createBlock using current NewestHash and Height
// and update NewestHash and Height for blockChain
func (b *blockChain) AddBlock() {
	block := createBlock(b.NewestHash, b.Height+1, getDifficulty(b))
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	persistBlockChain(b)
}

//------------ function for blockChain ------------------

// BlockChain get blockChain
// This function is for singleton pattern of blockChain
func BlockChain() *blockChain {
	once.Do(func() {
		b = &blockChain{
			Height: 0,
		}
		checkPoint := db.CheckPoint()
		if checkPoint == nil {
			b.AddBlock()
		} else {
			b.restoreFromBytes(checkPoint)
		}
	})
	return b
}

func persistBlockChain(b *blockChain) {
	db.SaveCheckPoint(utils.ToBytes(b))
}

//Blocks return all pointer of Blocks from DB
func Blocks(b *blockChain) []*Block {
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

func recalculateDifficulty(b *blockChain) int {
	allBlocks := Blocks(b)
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

func getDifficulty(b *blockChain) int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		return recalculateDifficulty(b)
	} else {
		return b.CurrentDifficulty
	}
}

func UTxOutsByAddress(address string, b *blockChain) []*UTxOut {
	var uTxOuts []*UTxOut
	creatorTxs := make(map[string]bool)
	for _, block := range Blocks(b) {
		for _, tx := range block.Transactions {
			for _, input := range tx.TxIns {
				if input.Owner == address {
					creatorTxs[input.TxID] = true
				}
			}

			for index, output := range tx.TxOuts {
				if output.Owner == address {
					if _, ok := creatorTxs[tx.ID]; !ok {
						uTxOut := &UTxOut{tx.ID, index, output.Amount}
						if !isOnMempool(uTxOut) {
							uTxOuts = append(uTxOuts, uTxOut)
						}
					}
				}
			}
		}
	}

	return uTxOuts
}

func BalanceByAddress(address string, b *blockChain) int {
	var amount int
	txOuts := UTxOutsByAddress(address, b)

	for _, txOut := range txOuts {
		amount += txOut.Amount
	}
	return amount
}
