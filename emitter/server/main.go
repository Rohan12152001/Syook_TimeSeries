package main

import (
	"flag"
	"fmt"
	"github.com/Rohan12152001/Syook_TimeSeries/emitter/utils"
	"github.com/gorilla/websocket"
	_ "html/template"
	"log"
	"net/http"
	"time"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options
var socketPool = map[*websocket.Conn]bool{}

type MessageStruct struct {
	enString string `json:"enString"`
}

func startEmitter(){
	ticker := time.NewTicker(time.Second*5)
	defer ticker.Stop()
	for{
		select {
		case <- ticker.C:
			// Form the encrypted string here
			encryptedString := utils.FormFinalString()
			fmt.Println("Emitter:", encryptedString)

			for obj := range socketPool{
				obj.WriteJSON(encryptedString)
			}
		}
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	// Add the connection to pool
	socketPool[ws] = true

	// helpful log statement to show connections
	log.Println("Client Connected")

	defer ws.Close()
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/ws", wsEndpoint)
	go startEmitter()
	log.Fatal(http.ListenAndServe(*addr, nil))
}