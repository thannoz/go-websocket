package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
)

// create a channel & a place to hold all connected users
var wsChan = make(chan WsPayload)
var clients = make(map[WebSocketConnection]string)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)

func checkError(err error) error {
	log.Fatal(err)
	return err
}

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Home(w http.ResponseWriter, request *http.Request) {
	err := renderPage(w, "home.jet", nil)
	if err != nil {
		checkError(err)
	}
}

type WebSocketConnection struct {
	*websocket.Conn
}

// WsJsonResponse defines the response sent back from websocket
type WsJsonResponse struct {
	Action      string `json:"action"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}

// WsJsonPayload define the information we will send to the server
type WsPayload struct {
	Action   string              `json:"action"`
	Username string              `json:"username"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"-"`
}

// upgradeConnection upgrade a normal http server connection to a websocket protocol
func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		checkError(err)
	}
	log.Println("Client connected to endpoint")

	var response WsJsonResponse
	response.Message = `<em><small>Connected to server</small></em>`

	// Added users to map when connected to ws-endpoint
	conn := WebSocketConnection{Conn: ws}
	clients[conn] = ""

	err = ws.WriteJSON(response)
	if err != nil {
		checkError(err)
	}

	go ListenforWs(&conn)
}

// take users away from WsEndpoint & put them into a goroutine
func ListenforWs(conn *WebSocketConnection) {
	// if connection (goroutine) dies, recover us
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	// listening for a payload
	var payload WsPayload
	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			// do nothing
		} else {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

func ListenToWsChannel() {
	var res WsJsonResponse

	// everytime we get a value from the channel, we use our res var to
	// populate it with some information
	for {
		event := <-wsChan
		res.Action = "Got here"
		res.Message = fmt.Sprintf("Some message, and action was %s", event.Action)
		broadcastToAll(res)
	}
}

// broadcast information to all users
func broadcastToAll(res WsJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(res)
		if err != nil {
			log.Println("websocket err")
			_ = client.Close()
			delete(clients, client)
		}
	}
}

func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		checkError(err)
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		checkError(err)
	}

	return nil
}
