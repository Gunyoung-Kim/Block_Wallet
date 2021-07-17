package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/Gunyoung-Kim/blockchain/db"
	"github.com/Gunyoung-Kim/blockchain/utils"
)

//ErrNotFound is error for not found
var ErrNotFound = errors.New("Not Found")

// Block is component of block chain
type Block struct {
	Height   int    `json:"height"`
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
}

//------------ receiver function for Block ------------------

func (b *Block) restoreFromBytes(data []byte) {
	utils.FromBytes(b, data)
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

// ----------- function for Block ----------------------------

//FindBlock find block from DB by Hash of Block
func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restoreFromBytes(blockBytes)
	return block, nil
}

//createBlock create Block using sha256
func createBlock(data string, prevHash string, height int) *Block {
	block := Block{
		Height:   height,
		Data:     data,
		Hash:     "",
		PrevHash: prevHash,
	}
	payload := block.Data + block.PrevHash + fmt.Sprint(block.Height)
	block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	block.persist()
	return &block
}
