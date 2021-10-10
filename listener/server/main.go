package main

import (
	"fmt"
	"github.com/Rohan12152001/Syook_TimeSeries/listener/manager/listener"
	"github.com/Rohan12152001/Syook_TimeSeries/listener/manager/listener/data"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var logger = logrus.New()
var upgrader = websocket.Upgrader{}
var socketPool = map[*websocket.Conn]bool{}
var listenerManager = listener.New()
var myListenerId = uuid.New().String()


// Here since we are using goroutines, so maybe we should use Atomic integers instead of normal integers,
// to ensure there is overwriting
var dataRecieved = 0
var dataUndecrypted = 0


type SucessRateStruct struct {
	DataRecieved int `json:"data_recieved"`
	DataUndecrypted int `json:"data_undecrypted"`
	SuccessRate int `json:"success_rate"`
}

func pushLiveDataIntoSockets(result data.LiveData){
	for client := range socketPool{
		client.WriteJSON(result)
	}
}

// Decrypt, Save & Emit
func decryptAndEmit(enStr string){
	splitArray := strings.Split(enStr, "|")

	for ind, objectString := range splitArray{
		if (ind+1) % 50 == 0{
			// Rate after every 50 messages
			successRate := 0
			if dataRecieved!=0{
				successRate = ((dataRecieved-dataUndecrypted)/dataRecieved)*100
			}
			logger.Infof("Cumulative Success Rate: %d", successRate)
		}

		decryptedObject, err := listenerManager.DecryptAndEmit(objectString, myListenerId)
		if err != nil {
			if err==listener.DiscardedError{
				dataUndecrypted+=1
				continue
			}
			logger.Error(err)
			return
		}
		dataRecieved += 1
		go pushLiveDataIntoSockets(decryptedObject)
	}
}

// Connect TO & Collect messages FROM emitter
func connectToEmitter(){

	host := os.Getenv("EMITTER_HOST")
	port := os.Getenv("EMITTER_PORT")

	//var addr = flag.String("addr", "host.docker.internal:8080", "http service address")

	addr := fmt.Sprintf("%v:%v", host,port)
	fmt.Println(">>>>> emitter", addr)

	u := url.URL{
		Scheme: "ws",
		Host: addr,
		Path: "/ws",
	}

	socketObject,_, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logger.Error(err)
	}

	defer socketObject.Close()

	for {
		_, message, err := socketObject.ReadMessage()
		// TODO: Handle this gracefully!
		if err != nil {
			logger.Error(err)
			return
		}
		fmt.Println("Got message")

		enStr := string(message)
		enStr = enStr[1 : len(enStr)-2]

		go decryptAndEmit(enStr)
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
			logger.Error(err)
			return
		}

		if err := conn.WriteMessage(messageType, message); err != nil {
			logger.Error(err)
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
		logger.Error("upgrade: ", err)
	}

	// Add the connection to pool
	socketPool[ws] = true

	// helpful log statement to show connections
	logger.Println("Client Connected")

	go reader(ws)
}

func homePage(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "./templates/index.html")
}

func setupRoutes() {
	http.HandleFunc("/ui", wsEndpoint)			// WebSocket for clients
	http.HandleFunc("/home", homePage)			// Client page
}

func main() {
	logger.Info("[START LISTENER SERVER]")

	setupRoutes()

	go connectToEmitter()

	// For free ports

	//freePort, err := freeport.GetFreePort()
	//if err != nil {
	//	log.Fatal(err)
	//}

	freePort := 53612

	portString := "0.0.0.0:" + strconv.Itoa(freePort)
	fmt.Println("Using port:", portString)
	logger.Fatal(http.ListenAndServe(portString, nil))
}