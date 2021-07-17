package blockchain

import (
	"errors"
	"strings"
	"time"

	"github.com/Gunyoung-Kim/blockchain/db"
	"github.com/Gunyoung-Kim/blockchain/utils"
)

//ErrNotFound is error for not found
var ErrNotFound = errors.New("Not Found")

// Block is component of block chain
type Block struct {
	Height     int    `json:"height"`
	Data       string `json:"data"`
	Hash       string `json:"hash"`
	PrevHash   string `json:"prevHash,omitempty"`
	Difficulty int    `json:"difficulty"`
	Nonce      int    `json:"nonce"`
	Timestamp  int    `json:"timestamp"`
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

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		b.Timestamp = int(time.Now().Unix())
		hash := utils.Hash(b)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}
	}
}

//createBlock create Block using sha256
func createBlock(data string, prevHash string, height int) *Block {
	block := Block{
		Height:     height,
		Data:       data,
		Hash:       "",
		PrevHash:   prevHash,
		Difficulty: BlockChain().difficulty(),
		Nonce:      0,
	}
	block.mine()
	block.persist()
	return &block
}
