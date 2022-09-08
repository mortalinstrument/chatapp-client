package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

type Message struct {
	Message   string
	Timestamp time.Time
	From      User
}

type User struct {
	Name      string
	LastLogin time.Time
	Active    bool
}

type Conn struct {
	sourceIp       net.Addr
	partner        *User
	connection     net.Conn
	messageChannel chan (Message)
}

//TODO for later on development
// var previousConnectionChannel = make(chan Conn, 1)

// var previousConnections = []Conn{}

// func addPreviousConnections() {
// 	for {
// 		v := <-previousConnectionChannel
// 		previousConnections = append(previousConnections, v)
// 	}
// }

// func checkForPreviousConnection(pC []Conn, c Conn) bool {
// 	for _, v := range pC {
// 		if v.partner.Name == c.partner.Name && v.sourceIp == c.sourceIp {
// 			return true
// 		}
// 	}
// 	return false
// }

func (c Conn) HandleRequest(log *os.File) error {
	// Make a buffer to hold incoming data. //TODO: BIGGER BYTE ARRAY
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	reqLen, err := c.connection.Read(buf)

	if err != nil {
		return err
	}
	input := buf[:reqLen]
	buffer := bytes.NewBuffer(input)
	decoder := gob.NewDecoder(buffer)

	recievedMessageObject := Message{}
	err = decoder.Decode(&recievedMessageObject)
	if err != nil {
		logger(fmt.Sprintf("couldn't decode message with err: %s", err), log)
		return err
	}
	logger(fmt.Sprintf("recieved message object with: %s:%s:%s", recievedMessageObject.From.Name, recievedMessageObject.Timestamp, recievedMessageObject.Message), log)
	c.connection.Close()
	return nil

	// //TODO: save connection and create a channel to save it in, or just add to existing channel out of list
	// wasConnectedBefore := checkForPreviousConnection(previousConnections, c)
	// if !wasConnectedBefore {
	// 	logger(fmt.Sprintf("client %s:%s hasnt connected before, creating new message channel...", c.sourceIp, c.partner.name), log)
	// 	channel := make(chan (message))
	// 	c.messageChannel = &channel
	// 	logger("done", log)
	// } else {
	// 	logger(fmt.Sprintf("client %s:%s has connected before, continuing", c.sourceIp, c.partner.name), log)
	// }

}

func (c Conn) SendRequest(msg string, log *os.File) error {
	logger(fmt.Sprintf("trying to send request with msg '%s' to %s:%s", msg, c.sourceIp, c.partner.Name), log)
	//TODO for later on development

	//lookup connection and if not found save in previousConnections
	// wasConnectedBefore := checkForPreviousConnection(previousConnections, c)
	// fmt.Println(wasConnectedBefore)
	// //TODO: create channel and add to connection before adding to previousConnections
	// if !wasConnectedBefore {
	// 	logger(fmt.Sprintf("client %s:%s hasnt been connected before, saving connection in previousConnection...", c.sourceIp, c.partner.Name), log)
	// 	previousConnectionChannel <- c
	// } else {
	// 	logger(fmt.Sprintf("client %s:%s has been connected before, lastLogin %s", c.sourceIp, c.partner.Name, c.partner.LastLogin), log)
	// }

	//check for emtpy message
	if msg == "" {
		return errors.New("Message was empty, exiting")
	}

	//create message object
	msgToSend := Message{
		Message:   msg,
		Timestamp: time.Now(),
		From:      *myself,
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(msgToSend); err != nil {
		logger(fmt.Sprintf("msg %s failed to encode with error: %s", msg, err), log)
	}

	//send to socket
	c.connection.Write(buf.Bytes())

	//read reply and print if there is any
	// message, err := bufio.NewReader(c.connection).ReadString('\n')
	// if err != nil {
	// 	logger(fmt.Sprintf("Error reading reply: %s", err.Error()), log)
	// }
	// if len(message) > 1 {
	// 	fmt.Print("-> " + message + "\n")
	// }

	c.connection.Close()
	return nil
}
