package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

func main() {
	// // Listen for incoming connections.
	// l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	// if err != nil {
	// 	fmt.Println("Error listening:", err.Error())
	// 	os.Exit(1)
	// }
	// // Close the listener when the application closes.
	// defer l.Close()
	// fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	// for {
	// 	// Listen for an incoming connection.
	// 	conn, err := l.Accept()
	// 	if err != nil {
	// 		fmt.Println("Error accepting: ", err.Error())
	// 		os.Exit(1)
	// 	}
	// 	// Handle connections in a new goroutine.
	// 	go handleRequest(conn)
	// }

	// self explanatory
	log := createLogFile()

	defer log.Close()

	//quickly connect to server to check for its availability
	checkForServer(log)
	//picking username
	userSetup(log)

	//loop for sending messages
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')

		sendRequest(text)
	}
}

func sendRequest(msg string) {
	//connect to server
	conn, err := net.Dial("tcp", "localhost:3333")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}

	//send to socket
	conn.Write([]byte(msg))

	//read reply and print if there is any
	message, _ := bufio.NewReader(conn).ReadString('\n')
	if len(message) > 1 {
		fmt.Print("-> " + message + "\n")
	}

	conn.Close()
}

func userSetup(log *os.File) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("choose a nickname: ")
	nick, _ := reader.ReadString('\n')

	conn, err := net.Dial("tcp", "localhost:3333")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}

	//send to socket
	conn.Write([]byte("0" + nick + " joined the chat"))

	reply, _ := bufio.NewReader(conn).ReadString('\n')
	//TODO: change if username is accepted by server
	if strings.Contains(reply, "Message recieved") {
		logger(fmt.Sprintf("nickname %s confirmed by server", nick), log)
		conn.Close()
	} else {
		logger(fmt.Sprintf("nickname %s couldn't be confirmed by server", nick), log)
	}

}

func checkForServer(log *os.File) {
	timeout := time.Duration(1 * time.Second)
	conn, err := net.DialTimeout(CONN_TYPE, CONN_HOST+":"+CONN_PORT, timeout)
	if err != nil {
		logger("Server was'nt available, Client shutting down", log)
		fmt.Println("Server not responding, exiting...")
		os.Exit(1)
	} else {
		logger(fmt.Sprintf("%s %s %s\n", CONN_HOST, "responding on port:", CONN_PORT), log)
		conn.Close()
	}
}

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
