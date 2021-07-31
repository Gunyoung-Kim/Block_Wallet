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

// Replace blocks of blockchain and reflect to DB
func (b *blockChain) Replace(blocks []*Block) {
	b.CurrentDifficulty = blocks[0].Difficulty
	b.Height = len(blocks)
	b.NewestHash = blocks[0].Hash
	persistBlockChain(b)

	for _, block := range blocks {
		block.persist()
	}
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

//persistBlockChain save checkpoint of blockchain to DB
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

//Transactions return all Transaxtion in blockChain
func Transactions(b *blockChain) []*Tx {
	var txs []*Tx
	for _, block := range Blocks(b) {
		txs = append(txs, block.Transactions...)
	}
	return txs
}

//FindTransaction return a transaction whose ID is corresponds with input targetID
func FindTransaction(b *blockChain, targetID string) *Tx {
	for _, tx := range Transactions(b) {
		if tx.ID == targetID {
			return tx
		}
	}

	return nil
}

//recalculateDifficulty recalculate difficulty of creating new block
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

//getDifficulty return current difficulty for creating new block.
//Return defaultDifficulty if there is no block yet.
//Return CurrentDifficulty of blockchain if it is not period of recalcualte difficulty
//Return new difficulty if it is period of recalculate difficulty
func getDifficulty(b *blockChain) int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		return recalculateDifficulty(b)
	} else {
		return b.CurrentDifficulty
	}
}

//UTxOutsByAddress return slice of UTxOut whose owner is given address
//exploring blockchain from newestblock to oldestblock if blocks transaction contain TxIn of address then add key,value to map
//if Tx contains txOut of address and its Id is not in map then add that txOut to UnusedTxOuts
func UTxOutsByAddress(address string, b *blockChain) []*UTxOut {
	var uTxOuts []*UTxOut
	creatorTxs := make(map[string]bool)
	for _, block := range Blocks(b) {
		for _, tx := range block.Transactions {
			for _, input := range tx.TxIns {
				if input.Signature == "COINBASE" {
					break
				}
				if FindTransaction(b, input.TxID).TxOuts[input.Index].Address == address {
					creatorTxs[input.TxID] = true
				}
			}

			for index, output := range tx.TxOuts {
				if output.Address == address {
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

//BalanceByAddress return balance of address which is calculated by slice of unused Txouts
func BalanceByAddress(address string, b *blockChain) int {
	var amount int
	txOuts := UTxOutsByAddress(address, b)

	for _, txOut := range txOuts {
		amount += txOut.Amount
	}
	return amount
}
