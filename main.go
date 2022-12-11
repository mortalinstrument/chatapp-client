package main

import (
	"bufio"
	"embed"
	"errors"
	"flag"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/seancfoley/ipaddress-go/ipaddr"
	"net"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

const (
	CONN_HOST    = "localhost"
	MESSAGE_PORT = "3333"
	CONN_TYPE    = "tcp"
)

type Config struct {
	Username string `envconfig:"CHAT_USERNAME"`
}

type Signaler struct {
	mu   sync.Mutex
	done bool
}

func (s *Signaler) setDoneTrue() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.done = true
}

func (s *Signaler) read() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.done
}

var emptyUserObject = User{}
var myself = &emptyUserObject

var addr = flag.String("addr", "localhost:7777", "service adress for frontendListener")
var directory = flag.String("d", "../chat/./dist", "the directory of static file to host")

var chat embed.FS

// creating new log file
var log = createLogFile()

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	//initalise channels for websocket communication on event with the frontend
	msgCannel := make(chan Message, 10)
	userChannel := make(chan User, 1)
	removeUserChannel := make(chan User, 1)

	var wg sync.WaitGroup
	signaler := Signaler{done: false}

	go func() {
		sig := <-sigs
		broadcast := calculateBroadcastAdress()
		conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{IP: broadcast, Port: 9999})
		if err != nil {
			logger(fmt.Sprintf("error while removing myself from explore: %s", err.Error()), log)
			os.Exit(1)
		}

		defer conn.Close()

		logger(fmt.Sprintf("sending removing explore to broadcast adress: %s", broadcast.String()), log)
		conn.Write([]byte(""))

		fmt.Print(sig)
		signaler.done = true
	}()

	var s Config
	err := envconfig.Process("chat", &s)
	if err != nil {
		fmt.Println("")
		os.Exit(1)
	}

	defer log.Close()

	//TODO: needs access to userchannel
	go explore(log)
	go listenForExplorers(log, userChannel, removeUserChannel)

	fmt.Println("username: " + s.Username)
	myself.Name = s.Username
	myself.Active = true
	myself.LastLogin = time.Now()
	myself.IP = getOwnIPAdress().String()

	//TODO: remove testing functions
	wg.Add(1) // Wait for 1 goroutine (thread) to be done before stopping to wait
	// start goroutine for messageListener
	go messageListener(&wg, &signaler, log, msgCannel)
	//instance new frontend and start listener from there
	go Frontend{msgCannel, userChannel, removeUserChannel}.frontendListener(log)

	//TODO: remove testing functions
	go func() {
		recipient := ipaddr.NewIPAddressString("127.0.0.1").GetAddress().GetNetIP()
		for i := 0; i < 1000; i++ {
			time.Sleep(6 * time.Second)
			sendRequest(fmt.Sprintf("Test Nachricht %i", i), &recipient, log)
		}
	}()

	//TODO: remove testing functions
	go func(signaler2 Signaler) {
		for !signaler2.read() {
			if len(exploredUsers) > 0 {
				fmt.Println("after if")
				recipient := ipaddr.NewIPAddressString(exploredUsers[0].IP).GetAddress().GetNetIP()
				fmt.Println("before for")
				for i := 0; i < 1000; i++ {
					fmt.Println(i)
					time.Sleep(6 * time.Second)
					sendRequest(fmt.Sprintf("Test Nachricht %i", i), &recipient, log)
				}
			}
		}
	}(signaler)

	wg.Wait()
}

func messageListener(wg *sync.WaitGroup, signaler *Signaler, log *os.File, msgChannel chan Message) error {
	time.Sleep(time.Duration(time.Second * 1))
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, "0.0.0.0:"+MESSAGE_PORT)
	if err != nil {
		fmt.Println("Error listening for messages: ", err.Error())
		return err
		wg.Done()
	}

	// Close the messageListener when the application closes.
	// defer functions run at the end of the parent-function
	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}
	}(l)

	logger(fmt.Sprintf("Listening for incoming messages on %s:%s", CONN_HOST, MESSAGE_PORT), log)
	for !signaler.read() {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting incoming message: ", err.Error())
			return err
			os.Exit(1)
		}
		connection := Conn{
			sourceIp:   conn.LocalAddr().String(),
			partner:    nil,
			connection: conn,
		}
		logger(fmt.Sprintf("handling request from %s", connection.sourceIp), log)
		// Handle connections in a new goroutine.
		go connection.HandleRequest(log, msgChannel)
	}

	wg.Done()
	return nil
}

func sendRequest(msg string, recipient *net.IP, log *os.File) (err error) {
	//connect to server
	conn, err := net.Dial(CONN_TYPE, recipient.String()+":"+MESSAGE_PORT)
	if err != nil {
		logger(fmt.Sprintf("Error connecting: %s", err.Error()), log)
		return errors.New("Error Connecting")
	}

	connection := Conn{
		sourceIp:   conn.LocalAddr().String(),
		partner:    myself,
		connection: conn,
	}
	connection.SendRequest(msg, log)
	return nil
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
// 	conn, err := net.Dial(CONN_TYPE, ":"+MESSAGE_PORT)
// 	if err != nil {
// 		logger("Server wasn't available, Client shutting down", log)
// 		fmt.Println("Server not responding, exiting...")
// 		conn.Close()
// 		os.Exit(1)
// 	} else {
// 		conn.Write([]byte("2"))
// 		logger(fmt.Sprintf("%s %s %s\n", CONN_HOST, "responding on port:", MESSAGE_PORT), log)
// 		conn.Close()
// 	}
// }

func createLogFile() *os.File {
	t := time.Now()
	os.Chdir("log")
	f, err := os.Create("log-" + t.Format("01-02-2006 15:04:05 Monday"))

	if err != nil {
		os.Exit(1)
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
