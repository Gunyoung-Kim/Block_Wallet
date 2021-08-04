package p2p

import (
	"encoding/json"
	"strings"

	"github.com/Gunyoung-Kim/blockchain/blockchain"
	"github.com/Gunyoung-Kim/blockchain/utils"
)

type MessageKind int

const (
	MessageNewestBlock MessageKind = iota
	MessageAllBlocksRequest
	MessageAllBlocksResponse
	MessageNewBlockNotify
	MessageNewTxNotify
	MessageNewPeerNotify
)

type Message struct {
	Kind    MessageKind
	Payload []byte
}

func makeMessage(kind MessageKind, payload interface{}) []byte {
	m := Message{
		Kind:    kind,
		Payload: utils.ToJSON(payload),
	}

	return utils.ToJSON(m)
}

func sendNewestBlock(p *peer) {
	b, err := blockchain.FindBlock(blockchain.BlockChain().NewestHash)
	utils.HandleError(err)

	m := makeMessage(MessageNewestBlock, b)

	p.inbox <- m
}

func requestAllBlocks(p *peer) {
	m := makeMessage(MessageAllBlocksRequest, nil)
	p.inbox <- m
}

func sendAllBlocks(p *peer) {
	m := makeMessage(MessageAllBlocksResponse, blockchain.Blocks(blockchain.BlockChain()))
	p.inbox <- m
}

func notifyNewBlock(b *blockchain.Block, p *peer) {
	m := makeMessage(MessageNewBlockNotify, b)
	p.inbox <- m
}

func notifyNewTx(tx *blockchain.Tx, p *peer) {
	m := makeMessage(MessageNewTxNotify, tx)
	p.inbox <- m
}

func notifyNewPeer(address string, p *peer) {
	m := makeMessage(MessageNewPeerNotify, address)
	p.inbox <- m
}

func handleMsg(m *Message, p *peer) {
	switch m.Kind {
	case MessageNewestBlock:
		var payload blockchain.Block
		utils.HandleError(json.Unmarshal(m.Payload, &payload))
		b, err := blockchain.FindBlock(blockchain.BlockChain().NewestHash)
		utils.HandleError(err)

		if payload.Height >= b.Height {
			requestAllBlocks(p)
		} else {
			sendNewestBlock(p)
		}
	case MessageAllBlocksRequest:
		sendAllBlocks(p)
	case MessageAllBlocksResponse:
		var payload []*blockchain.Block
		json.Unmarshal(m.Payload, &payload)
		utils.HandleError(json.Unmarshal(m.Payload, &payload))
		blockchain.BlockChain().Replace(payload)
	case MessageNewBlockNotify:
		var payload *blockchain.Block
		utils.HandleError(json.Unmarshal(m.Payload, &payload))
		blockchain.BlockChain().AddPeerBlock(payload)
	case MessageNewTxNotify:
		var payload *blockchain.Tx
		utils.HandleError(json.Unmarshal(m.Payload, &payload))
		blockchain.Mempool().AddPeerTx(payload)
	case MessageNewPeerNotify:
		var payload string
		utils.HandleError(json.Unmarshal(m.Payload, &payload))
		parts := strings.Split(payload, ":")
		AddPeer(parts[0], parts[1], parts[2], false)
	}

}
