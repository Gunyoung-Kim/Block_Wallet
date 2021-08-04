package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Gunyoung-Kim/blockchain/blockchain"
	"github.com/Gunyoung-Kim/blockchain/p2p"
	"github.com/Gunyoung-Kim/blockchain/utils"
	"github.com/Gunyoung-Kim/blockchain/wallet"
	"github.com/gorilla/mux"
)

var port string

type url string

func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

// urlDescription is reponse entity for url Description
type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

// balanceResponse is response entity for balance
type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

// myWalletResponse is reponse entity for wallet status
type myWalletResponse struct {
	Address string `json:"address"`
}

// errorResponse is reponse entity for error message
type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type addTxPayload struct {
	To     string
	Amount int
}

type addPeerPayLoad struct {
	Address string `json:"address"`
	Port    string `json:"port"`
}

// documentation show all url possible and its description in this api server
func documentation(rw http.ResponseWriter, req *http.Request) {
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         url("/status"),
			Method:      "GET",
			Description: "See the status of blockchain",
		},
		{
			URL:         url("/blocks"),
			Method:      "GET",
			Description: "See All Blocks",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add a Block",
			Payload:     "data:string",
		},
		{
			URL:         url("/blocks/{hash}"),
			Method:      "GET",
			Description: "See a Block",
		},
		{
			URL:         url("/balance/{address}"),
			Method:      "GET",
			Description: "Get TxOuts for an address",
		},
		{
			URL:         url("/ws"),
			Method:      "GET",
			Description: "Upgrade to web sockets",
		},
	}

	json.NewEncoder(rw).Encode(data)
}

// status return status of current blockchain in server
func status(rw http.ResponseWriter, req *http.Request) {
	blockchain.Status(blockchain.BlockChain(), rw)
}

// blocks take two methods
// if request's method is GET, then return all blocks in blockChain
// if request's method is POST, then add new block to blockChain
func blocks(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		json.NewEncoder(rw).Encode(blockchain.Blocks(blockchain.BlockChain()))
	case "POST":
		newBlock := blockchain.BlockChain().AddBlock()
		p2p.BroadcastNewBlock(newBlock)
		rw.WriteHeader(http.StatusCreated)
	}
}

// block return status of a block by hash
// it returns {@code blockChain.ErrNotFound} if there is no such block
func block(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	hash := vars["hash"]
	encoder := json.NewEncoder(rw)
	block, err := blockchain.FindBlock(hash)
	if err == blockchain.ErrNotFound {
		encoder.Encode(errorResponse{ErrorMessage: fmt.Sprint(err)})
	} else {
		encoder.Encode(block)
	}
}

// balance return current balance of address
// if request query contains total, then it returns amount of balance
// if it doesn't contain, then return list of unused transaction output
func balance(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	address := vars["address"]
	isTotal := req.URL.Query().Get("total")

	if isTotal == "true" {
		amount := blockchain.BalanceByAddress(address, blockchain.BlockChain())
		balanceRes := balanceResponse{Address: address, Balance: amount}
		utils.HandleError(json.NewEncoder(rw).Encode(balanceRes))
	} else {
		utils.HandleError(json.NewEncoder(rw).Encode(blockchain.UTxOutsByAddress(address, blockchain.BlockChain())))
	}
}

// mempool return all transactions in Mempool
func mempool(rw http.ResponseWriter, req *http.Request) {
	utils.HandleError(json.NewEncoder(rw).Encode(blockchain.Mempool().Txs))
}

// transactions add new transaction in Mempool
// it return status created
// if there comes error while creaing transaction, then it return errorMsg with status BadRequest
func transactions(rw http.ResponseWriter, req *http.Request) {
	var payload addTxPayload
	utils.HandleError(json.NewDecoder(req.Body).Decode(&payload))
	tx, err := blockchain.Mempool().AddTx(payload.To, payload.Amount)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(errorResponse{err.Error()})
		return
	}

	p2p.BroadcastNewTx(tx)

	rw.WriteHeader(http.StatusCreated)
}

// myWallet return address of wallet which is made by public key
func myWallet(rw http.ResponseWriter, req *http.Request) {
	address := wallet.Wallet().Address
	json.NewEncoder(rw).Encode(myWalletResponse{Address: address})
}

func peers(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		var payload addPeerPayLoad
		json.NewDecoder(req.Body).Decode(&payload)
		p2p.AddPeer(payload.Address, payload.Port, port[1:], true)
		rw.WriteHeader(http.StatusOK)
	case "GET":
		json.NewEncoder(rw).Encode(p2p.AllPeers(&p2p.Peers))
	}
}

// jsonContentTypeMiddleWare define content type of all response to  {@code application/json}
// this function use 'Adapter Pattern'
func jsonContentTypeMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, req)
	})
}

// loggerMiddleWare log request URL
func loggerMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fmt.Println(req.URL)
		next.ServeHTTP(rw, req)
	})
}

//Start REST API
func Start(portNum int) {
	router := mux.NewRouter()
	port = fmt.Sprintf(":%d", portNum)
	router.Use(jsonContentTypeMiddleWare, loggerMiddleWare)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", status).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/balance/{address}", balance).Methods("GET")
	router.HandleFunc("/mempool", mempool).Methods("GET")
	router.HandleFunc("/wallet", myWallet).Methods("GET")
	router.HandleFunc("/transactions", transactions).Methods("POST")
	router.HandleFunc("/ws", p2p.Upgrade).Methods("GET")
	router.HandleFunc("/peers", peers).Methods("GET", "POST")
	fmt.Printf("REST Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
