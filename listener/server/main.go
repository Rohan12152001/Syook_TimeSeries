package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/phayes/freeport"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

var logger = logrus.New()
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var socketPool = map[*websocket.Conn]bool{}

/*
Both steps need to be at service layer, and return the object here at ENDPOINT layer & use it
// 1. routine for validation

// 2. Adding messages into DB

// 3. sending decrypted messages to all UI's
*/

// Connect TO & Collect messages FROM emitter
func connectToEmitter(){
	var addr = flag.String("addr", "localhost:8080", "http service address")

	u := url.URL{
		Scheme: "ws",
		Host: *addr,
		Path: "/ws",
	}

	socketObject,_, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
	}

	defer socketObject.Close()

	for {
		_, message, err := socketObject.ReadMessage()
		if err != nil {
			log.Println("error: ", err)
			return
		}
		fmt.Println("message: ", message)
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
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

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
	http.HandleFunc("/ui", wsEndpoint)			// WS for clients
}

func main() {
	logger.Info("[START LISTENER SERVER]")

	setupRoutes()

	// TODO: connectEmitter
	go connectToEmitter()

	freePort, err := freeport.GetFreePort()
	if err != nil {
		log.Fatal(err)
	}

	portString := ":" + strconv.Itoa(freePort)
	fmt.Println("Using port:", portString)
	logger.Fatal(http.ListenAndServe(portString, nil))
}