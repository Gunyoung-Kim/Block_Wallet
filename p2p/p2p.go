package p2p

import (
	"fmt"
	"net/http"

	"github.com/Gunyoung-Kim/blockchain/utils"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

//Upgrade turn http/https connection into web socket connection
func Upgrade(rw http.ResponseWriter, req *http.Request) {
	// Port :3000 will upgrade the request from :4000
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(rw, req, nil)
	utils.HandleError(err)
	initPeer(conn, "", "")
}

func AddPeer(address, port string) {
	// Port :4000 is request an upgrade from the port :3000
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws", address, port), nil)
	utils.HandleError(err)
	initPeer(conn, address, port)
}
