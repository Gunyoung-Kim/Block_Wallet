package p2p

import (
	"net/http"

	"github.com/Gunyoung-Kim/blockchain/utils"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

//Upgrade turn http/https connection into web socket connection
func Upgrade(rw http.ResponseWriter, req *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	_, err := upgrader.Upgrade(rw, req, nil)
	utils.HandleError(err)
}
