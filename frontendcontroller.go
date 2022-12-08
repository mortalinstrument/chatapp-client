package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/seancfoley/ipaddress-go/ipaddr"
	"net/http"
	"os"
	"sync"
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
	messageChannel chan Message
}

func (f Frontend) frontendListener(wg sync.WaitGroup, log *os.File) {
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(*directory)))
	http.HandleFunc("/c", f.serveWs)

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		logger("error while trying to listen for frontend", log)
		wg.Done()
	}
	logger(fmt.Sprintf("Listening for Frontend Delivery on %s", addr), log)

}

func readPump(conn *websocket.Conn) {
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

func (f Frontend) writePump(conn *websocket.Conn) {
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

func (frontend Frontend) serveWs(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go frontend.writePump(conn)
	go readPump(conn)

}
