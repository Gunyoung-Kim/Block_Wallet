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

type addBlockBody struct {
	Message string
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

func documentation(rw http.ResponseWriter, req *http.Request) {
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See Documentation",
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
			URL:         url("blocks/{hash}"),
			Method:      "GET",
			Description: "See a Block",
		},
	}

	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		json.NewEncoder(rw).Encode(blockchain.BlockChain().Blocks())
	case "POST":
		var addBlockBody addBlockBody
		utils.HandleError(json.NewDecoder(req.Body).Decode(&addBlockBody))
		blockchain.BlockChain().AddBlock(addBlockBody.Message)
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
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	fmt.Printf("REST Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
