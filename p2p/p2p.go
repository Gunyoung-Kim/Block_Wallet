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
	openPort := req.URL.Query().Get("openPort")
	ip := utils.Splitter(req.RemoteAddr, ":", 0)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return openPort != "" && ip != ""
	}
	conn, err := upgrader.Upgrade(rw, req, nil)
	utils.HandleError(err)
	initPeer(conn, ip, openPort)
}

func AddPeer(address, port, openPort string) {
	// Port :4000 is request an upgrade from the port :3000
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort[1:]), nil)
	utils.HandleError(err)
	p := initPeer(conn, address, port)
	sendNewestBlock(p)
}
