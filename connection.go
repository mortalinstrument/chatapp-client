package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
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

type IncomingMessage struct {
	Message   string
	Timestamp string
	ToIP      string
	ToName    string
}

type User struct {
	Name      string
	IP        string
	LastLogin time.Time
	Active    bool
	Messages  []Message
}

type Conn struct {
	sourceIp   string
	partner    *User
	connection net.Conn
}

func (c Conn) HandleRequest(log *os.File, msgChannel chan Message) error {
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
	msgChannel <- recievedMessageObject
	logger(fmt.Sprintf("recieved message object with: %s:%s:%s", recievedMessageObject.From.Name, recievedMessageObject.Timestamp, recievedMessageObject.Message), log)
	c.connection.Close()
	return nil
}

func (c Conn) SendRequest(msg string, log *os.File) error {
	logger(fmt.Sprintf("trying to send request with msg '%s' to %s:%s", msg, c.sourceIp, c.partner.Name), log)

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
	enc := json.NewEncoder(&buf)

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
