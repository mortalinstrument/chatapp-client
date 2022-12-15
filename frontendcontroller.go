package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/seancfoley/ipaddress-go/ipaddr"
	"net/http"
	"os"
	"time"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 102400
)

type Frontend struct {
	messageChannel    chan Message
	userChannel       chan User
	removeUserChannel chan User
}

func (f Frontend) frontendListener(log *os.File) {
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(*directory)))
	http.HandleFunc("/whoami", f.serveUserInfo)
	http.HandleFunc("/whothere", f.serveAllUsers)
	http.HandleFunc("/c", f.serveMessagesWs)
	http.HandleFunc("/cu", f.serveUsersWs)

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		logger("error while trying to listen for frontend", log)
	}
	logger(fmt.Sprintf("Listening for Frontend Delivery on %s", addr), log)

}

func readMessagePump(conn *websocket.Conn) {
	defer conn.Close()
	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, reader, err := conn.NextReader()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println(err.Error())
			}
			fmt.Println(err.Error())
			return
		}
		fmt.Println(reader)
		decoder := json.NewDecoder(reader)
		receivedMessageObject := IncomingMessage{}
		err = decoder.Decode(&receivedMessageObject) //decodes Request-Body into Message
		if err != nil {
			fmt.Println("cannot Decode Message")
			fmt.Println(receivedMessageObject)
			return
		}

		fmt.Println(receivedMessageObject)

		recipientIp := ipaddr.NewIPAddressString(receivedMessageObject.ToIP).GetAddress().GetNetIP()
		sendRequest(receivedMessageObject.Message, &recipientIp, log)
	}
}

func (f Frontend) writeMessagePump(conn *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case msg := <-f.messageChannel:
			conn.SetWriteDeadline(time.Now().Add(writeWait))

			w, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			// Encoding Message struct to JSON-Object
			data, _ := json.Marshal(msg)

			// write JSON-Object to frontend
			_, err = w.Write(data)
			if err != nil {
				fmt.Println("Writing error to frontend")
				os.Exit(1)
			}

			fmt.Println(string(data))

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		}
	}
}

// TODO:FIX WEBSOCKET FOR USERS
func (f Frontend) writeUserPump(conn *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		select {
		case removeUser := <-f.removeUserChannel:
			conn.SetWriteDeadline(time.Now().Add(writeWait))

			w, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			_, err = w.Write([]byte("remove user with ip:" + removeUser.IP))
			if err != nil {
				fmt.Println("Writing error to frontend")
				os.Exit(1)
			}

			fmt.Println("written userremoval to frontend")

			if err := w.Close(); err != nil {
				return
			}
		case user := <-f.userChannel:
			conn.SetWriteDeadline(time.Now().Add(writeWait))

			w, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			// Encoding Message struct to JSON-Object
			data, _ := json.Marshal(user)

			// write JSON-Object to frontend
			_, err = w.Write(data)
			if err != nil {
				fmt.Println("Writing error to frontend")
				os.Exit(1)
			}

			fmt.Println(string(data))

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (frontend Frontend) serveUserInfo(w http.ResponseWriter, r *http.Request) {
	userInfo, err := json.Marshal(*myself)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
	w.Write(userInfo)
}

func (frontend Frontend) serveMessagesWs(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go frontend.writeMessagePump(conn)
	go readMessagePump(conn)
}

func (frontend Frontend) serveUsersWs(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	go frontend.writeUserPump(conn)
}

func (frontend Frontend) serveAllUsers(w http.ResponseWriter, r *http.Request) {
	allUsers, err := json.Marshal(exploredUsers)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
	w.Write(allUsers)
}
