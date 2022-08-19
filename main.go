package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

var myself *User

func main() {
	// self explanatory
	log := createLogFile()

	defer log.Close()
	// start goroutine for listener
	go listener(log)

	//picking username
	userSetup(log)

	sendRequest("test", log)

	// //loop for sending messages
	// for {
	// 	reader := bufio.NewReader(os.Stdin)
	// 	fmt.Print(">> ")
	// 	text, _ := reader.ReadString('\n')

	// 	sendRequest(text)
	// }
}

func listener(log *os.File) {
	time.Sleep(time.Duration(time.Second * 2))
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	//fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		connection := Conn{
			sourceIp:   conn.RemoteAddr().Network(),
			partner:    nil,
			connection: conn,
			channel:    nil,
		}

		go connection.HandleRequest()
	}
}

func sendRequest(msg string, log *os.File) {
	//connect to server
	conn, err := net.Dial(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		logger(fmt.Sprintf("Error connecting: %s", err.Error()), log)
		os.Exit(1)
	}

	connection := Conn{
		sourceIp:   conn.LocalAddr().Network(),
		partner:    myself,
		connection: conn,
		channel:    nil,
	}

	connection.SendRequest(msg, log)
}

func contactUser() {
	// func for first contact, channel for messages should be created here
}

func userSetup(log *os.File) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("choose a nickname: ")
	nick, _ := reader.ReadString('\n')
	logger(fmt.Sprintf("nickname set as %s", nick), log)

	myself.name = nick
	myself.lastLogin = time.Now()
	myself.active = true
}

// func checkForServer(log *os.File) {
// 	conn, err := net.Dial(CONN_TYPE, ":"+CONN_PORT)
// 	if err != nil {
// 		logger("Server was'nt available, Client shutting down", log)
// 		fmt.Println("Server not responding, exiting...")
// 		conn.Close()
// 		os.Exit(1)
// 	} else {
// 		conn.Write([]byte("2"))
// 		logger(fmt.Sprintf("%s %s %s\n", CONN_HOST, "responding on port:", CONN_PORT), log)
// 		conn.Close()
// 	}
// }

func createLogFile() *os.File {
	t := time.Now()
	f, err := os.Create("log-" + t.Format("01-02-2006 15:04:05 Monday"))

	if err != nil {
		log.Fatal(err)
	}

	return f
}

func logger(msg string, log *os.File) {
	h, m, s := time.Now().Hour(), time.Now().Minute(), time.Now().Second()
	log.WriteString(fmt.Sprintf("%d:%d:%d : %s \n", h, m, s, msg))
}

// // Handles incoming requests.
// func handleRequest(conn net.Conn) {
// 	// Make a buffer to hold incoming data.
// 	buf := make([]byte, 1024)
// 	// Read the incoming connection into the buffer.
// 	reqLen, err := conn.Read(buf)
// 	if err != nil {
// 		fmt.Println("Error reading:", err.Error())
// 	}
// 	// Send a response back to person contacting us.
// 	conn.Write([]byte(fmt.Sprintf("Message received. Length: %d", reqLen)))
// 	// Close the connection when you're done with it.
// 	conn.Close()
// }
