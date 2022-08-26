package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"time"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

var emptyUserObject = User{}
var myself = &emptyUserObject

func main() {
	// self explanatory
	log := createLogFile()

	defer log.Close()
	// start goroutine for listener
	go listener(log)
	time.Sleep(time.Duration(time.Second * 3))

	//picking username
	userSetup(log)

	//go addPreviousConnections()

	//loop for sending messages
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')

		go sendRequest(text, log)
	}
}

func listener(log *os.File) {
	time.Sleep(time.Duration(time.Second * 1))
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	logger(fmt.Sprintf("Listening for incoming connections on %s:%s", CONN_HOST, CONN_PORT), log)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		//TODO: evtl partner mitsenden, um hier eintragen zu können, danach nachricht
		connection := Conn{
			sourceIp:   conn.RemoteAddr(),
			partner:    nil,
			connection: conn,
		}
		//logger(fmt.Sprintf("handling request from %s sent by user with name %s", connection.sourceIp, connection.partner.name), log)
		// Handle connections in a new goroutine.
		go connection.HandleRequest(log)
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
		sourceIp:   conn.LocalAddr(),
		partner:    myself,
		connection: conn,
	}

	connection.SendRequest(msg, log)
}

func firstContact(conn Conn) {
	// func for first contact, channel for messages should be created here

}

func userSetup(log *os.File) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("choose a nickname: ")
	nick, _ := reader.ReadString('\n')
	logger(fmt.Sprintf("nickname set as %s", nick), log)

	myself.Name = nick
	myself.LastLogin = time.Now()
	myself.Active = true
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
	_, file, line, _ := runtime.Caller(1)
	log.WriteString(fmt.Sprintf("%s Zeile:%d  	%d:%d:%d : %s \n", file, line, h, m, s, msg))

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
