package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Gunyoung-Kim/blockchain/blockchain"
	"github.com/Gunyoung-Kim/blockchain/utils"
	"github.com/gorilla/mux"
)

var port string

type url string

func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type addTxPayload struct {
	To     string
	Amount int
}

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
	}

	json.NewEncoder(rw).Encode(data)
}

func status(rw http.ResponseWriter, req *http.Request) {
	json.NewEncoder(rw).Encode(blockchain.BlockChain())
}

func blocks(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		json.NewEncoder(rw).Encode(blockchain.Blocks(blockchain.BlockChain()))
	case "POST":
		blockchain.BlockChain().AddBlock()
		rw.WriteHeader(http.StatusCreated)
	}
}

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

func mempool(rw http.ResponseWriter, req *http.Request) {
	utils.HandleError(json.NewEncoder(rw).Encode(blockchain.Mempool.Txs))
}

func transactions(rw http.ResponseWriter, req *http.Request) {
	var payload addTxPayload
	utils.HandleError(json.NewDecoder(req.Body).Decode(&payload))
	err := blockchain.Mempool.AddTx(payload.To, payload.Amount)
	if err != nil {
		json.NewEncoder(rw).Encode(errorResponse{"not enough funds"})
	}

	rw.WriteHeader(http.StatusCreated)
}

// Adapter Pattern !!!
func jsonContentTypeMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, req)
	})
}

//Start REST API
func Start(portNum int) {
	router := mux.NewRouter()
	port = fmt.Sprintf(":%d", portNum)
	router.Use(jsonContentTypeMiddleWare)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", status).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/balance/{address}", balance)
	router.HandleFunc("/mempool", mempool)
	router.HandleFunc("/transactions", transactions).Methods("POST")
	fmt.Printf("REST Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
