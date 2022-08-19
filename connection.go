package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

type message struct {
	message   string
	timestamp time.Time
}

type User struct {
	name      string
	lastLogin time.Time
	active    bool
}

type Conn struct {
	sourceIp   string
	partner    *User
	connection net.Conn
	channel    chan (message)
}

var previousConnections = []Conn{}

func (c Conn) HandleRequest() {

	// Make a buffer to hold incoming data. //TODO: BIGGER BYTE ARRAY
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	reqLen, err := c.connection.Read(buf)

	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	jobCode := string(buf[:1])

	if jobCode == "0" { // UserSetup
		//remove -2 if not sent via console (\n)
		content := string(buf[1 : reqLen-2])
		c.partner = &User{
			name:      content,
			lastLogin: time.Now(),
			active:    true,
		}

		previousConnections = append(previousConnections, c)
		c.connection.Write([]byte("UserSetup successful"))
	} else if jobCode == "1" { //Chat Message
		// remove -2 if not sent via console (\n)
		content := string(buf[:reqLen-2])
		// send a response back to person contacting us.
		c.connection.Write([]byte(fmt.Sprintf("Message received. Length: %d", reqLen)))
		fmt.Printf("Message contents: %q\n", content)
		c.connection.Close()
	} else if jobCode == "2" { //ping
		c.connection.Write([]byte("ChatService is ready"))
		fmt.Printf("%s started client and successfully test connected\n", c.connection.RemoteAddr())
		c.connection.Close()
	}
}

func (c Conn) SendRequest(msg string, log *os.File) {
	//send to socket
	c.connection.Write([]byte(msg))

	//read reply and print if there is any
	message, err := bufio.NewReader(c.connection).ReadString('\n')
	if err != nil {
		logger(fmt.Sprintf("Error reading reply: %s", err.Error()), log)
	}
	if len(message) > 1 {
		fmt.Print("-> " + message + "\n")
	}

	c.connection.Close()
}
