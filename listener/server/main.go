package main

import (
	"flag"
	"fmt"
	"github.com/Rohan12152001/Syook_TimeSeries/listener/manager/listener"
	"github.com/Rohan12152001/Syook_TimeSeries/listener/manager/listener/data"
	"github.com/gorilla/websocket"
	"github.com/phayes/freeport"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var logger = logrus.New()
var upgrader = websocket.Upgrader{}
var socketPool = map[*websocket.Conn]bool{}

type MessageStruct struct {
	enString string `json:"enString"`
}


// Decrypt, Save & Emit
func decryptAndEmit(enStr string){
	splitArray := strings.Split(enStr, "|")

	//TODO: (ASK should we pass a channel to the manager
	// & while loop on the channel below here ?)

	liveChannel := make(chan data.LiveData)
	go listener.DecryptAndEmit(splitArray, &liveChannel)		// On service layer Todo: Looks fine ?

	for result := range liveChannel {
		fmt.Println("RESULT: ", result)
		// TODO: Here iterate socket pool and emit to all
	}
	fmt.Println("DONE!")
}

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
			log.Println(err)
			return
		}

		enStr := string(message)
		enStr = enStr[1 : len(enStr)-2]

		decryptAndEmit(enStr)
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
	http.HandleFunc("/ui", wsEndpoint)			// WebSocket for clients
}

func main() {
	logger.Info("[START LISTENER SERVER]")

	setupRoutes()

	go connectToEmitter()

	freePort, err := freeport.GetFreePort()
	if err != nil {
		log.Fatal(err)
	}

	portString := ":" + strconv.Itoa(freePort)
	fmt.Println("Using port:", portString)
	logger.Fatal(http.ListenAndServe(portString, nil))
}