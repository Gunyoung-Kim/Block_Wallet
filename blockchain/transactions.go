package blockchain

import (
	"errors"
	"time"

	"github.com/Gunyoung-Kim/blockchain/utils"
)

const (
	minerReward int = 50
)

type mempool struct {
	Txs []*Tx
}

//Mempool slice of Tx which is not confirmed
var Mempool *mempool = &mempool{}

//AddTx add new transaction to mempool
func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx("kim", to, amount)

	if err != nil {
		return err
	}

	m.Txs = append(m.Txs, tx)
	return nil
}

//txToConfirm confirm all transactions in mempool
//get all transaction from mempool and add coinbaseTx then return transactions,
//then initialize mempool
func (m *mempool) txToConfirm() []*Tx {
	coinbase := makeCoinbaseTx("kim")
	txs := m.Txs
	txs = append(txs, coinbase)
	m.Txs = nil
	return txs
}

//Tx is transaction
type Tx struct {
	ID        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

//getID create ID for Tx by hashing another field of Tx
func (t *Tx) getID() {
	t.ID = utils.Hash(t)
}

//TxIn represents input for transaction
type TxIn struct {
	TxID  string `json:"txID"`
	Index int    `json:"index"`
	Owner string `json:"owner"`
}

//TxOut represents output for transaction
type TxOut struct {
	Owner  string
	Amount int
}

//UTxOut represents TxOut which is not used for input of transaction
type UTxOut struct {
	TxID   string
	Index  int
	Amount int
}

//isOnMempool check UTxOut is in TxIns in Tx in mempool before add to result of unusedTxOut
func isOnMempool(uTxOut *UTxOut) bool {
	for _, tx := range Mempool.Txs {
		for _, input := range tx.TxIns {
			if input.TxID == uTxOut.TxID && input.Index == uTxOut.Index {
				return true
			}
		}
	}

	return false
}

//makeCoinbaseTx make Tx from coinbase for miner
func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"", -1, "COINBASE"},
	}

	txOuts := []*TxOut{
		{address, minerReward},
	}

	tx := Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}

	tx.getID()
	return &tx
}

//makeTx make transction for input amount
//first check from has enough balance by blockchain
//then get all unusedTxOuts and add one to one, make txIn until total is bigger than or equal to amount
//if total is bigger than amount then append changeTxOut to txOuts of new Tx
func makeTx(from, to string, amount int) (*Tx, error) {
	if BalanceByAddress(from, BlockChain()) < amount {
		return nil, errors.New("Not enough balance")
	}

	var txOuts []*TxOut
	var txIns []*TxIn
	total := 0
	uTxOuts := UTxOutsByAddress(from, BlockChain())
	for _, uTxOut := range uTxOuts {
		if total >= amount {
			break
		}
		txIn := &TxIn{uTxOut.TxID, uTxOut.Index, from}
		txIns = append(txIns, txIn)
		total += uTxOut.Amount
	}

	if change := total - amount; change != 0 {
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}

	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)
	tx := &Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getID()
	return tx, nil
}
