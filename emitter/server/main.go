package main

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"time"
)

var logger = logrus.New()
var upgrader = websocket.Upgrader{}
var socketPool = map[*websocket.Conn]bool{}

func startEmitter(){
	// Form the encrypted string here
	//encryptedString := utils.FormFinalString()	// IMP

	// emit in every 10 secs to all the clientSockets
	for{
		time.Sleep(10)
		//for obj := range socketPool{
		//	// TODO: Here
		//}
	}

	// tmp, emit hello after every 10secs
	for{
		time.Sleep(3)
		for obj := range socketPool{
			obj.WriteJSON("hello")
		}
	}
}

func reader(conn *websocket.Conn) {
	defer func() {
		delete(socketPool, conn)
		conn.Close()
	}()
	for {
		// read in a message
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Println(err)
			return
		}
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	//upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	// Add the connection to pool
	socketPool[ws] = true

	// helpful log statement to show connections
	log.Println("Client Connected")

	go reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	logger.Info("[START EMITTER SERVER]")
	setupRoutes()
	startEmitter()
	logger.Fatal(http.ListenAndServe(":8080", nil))
}







