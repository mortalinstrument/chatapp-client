package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/seancfoley/ipaddress-go/ipaddr"
	"net"
	"os"
	"strings"
	"time"
)

var exploredUsers []User

func explore(log *os.File) {
	broadcast := calculateBroadcastAdress()
	fmt.Println(broadcast)

	for {
		time.Sleep(time.Second * 30)

		conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{IP: broadcast, Port: 9999})
		if err != nil {
			logger(fmt.Sprintf("error while exploring: %s", err.Error()), log)
			os.Exit(1)
		}

		defer conn.Close()

		logger(fmt.Sprintf("sending explore to broadcast adress: %s", broadcast.String()), log)
		conn.Write([]byte(myself.Name))
	}
}

func listenForExplorers(log *os.File) {
	pc, err := net.ListenPacket("udp4", ":9999")
	logger("listening for other clients to explore on port 9999", log)
	if err != nil {
		fmt.Println("Error while listening for other clients: " + err.Error())
	}
	defer pc.Close()

	for {
		data := make([]byte, 4096)
		read, remoteAddr, err := pc.ReadFrom(data)
		if err != nil {
			os.Exit(1)
		}
		ipWithoutPort := strings.Split(remoteAddr.String(), ":")[0]
		newUser := User{
			Name:      string(data[:read]),
			IP:        ipWithoutPort,
			LastLogin: time.Now(),
			Active:    true,
		}
		userDoesAlreadyExist := func() bool {
			for _, u := range exploredUsers {
				if u.Name == newUser.Name && u.IP == newUser.IP {
					return true
					break
				}
			}
			return false
		}
		if !(newUser.Name == myself.Name) || !(newUser.IP == myself.IP) {
			if !userDoesAlreadyExist() {
				exploredUsers = append(exploredUsers, newUser)
				logger("got exploring message from another client, creating user: "+newUser.Name, log)
			} else {
				logger("recieved Broadcast, but ignored it because user already exists", log)
			}
		} else {
			logger("recieved Broadcast, but ignored it because it is myself", log)
		}
	}
}

func getOwnIPAdress() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

func calculateBroadcastAdress() net.IP {
	ip := getOwnIPAdress()

	interfaces, err := net.Interfaces()
	var localIps []net.Addr
	if err != nil {
		os.Exit(1)
	}

	for _, iface := range interfaces {
		addrs, _ := iface.Addrs()
		localIps = append(localIps, addrs...)
	}

	for _, localIp := range localIps {
		cidrStr := localIp.String()
		maskAddr := ipaddr.NewIPAddressString(cidrStr).GetAddress().GetNetworkMask() // add .GetNetIP() for net.IP

		println("mask: " + maskAddr.String() + "ip: " + ipaddr.NewIPAddressString(cidrStr).GetAddress().String())

		_, ipnet, err := net.ParseCIDR(localIp.String())
		if err != nil {
			os.Exit(1)
		}

		if ipnet.Contains(ip) {
			println("mask: " + maskAddr.String() + "ip: " + ip.String())

			lastAddr, err := lastAddr(ipnet)
			if err != nil {
				os.Exit(1)
			}
			return lastAddr
		}
	}
	return net.IP{}
}

func lastAddr(n *net.IPNet) (net.IP, error) { // works when the n is a prefix, otherwise...
	if n.IP.To4() == nil {
		return net.IP{}, errors.New("does not support IPv6 addresses.")
	}
	ip := make(net.IP, len(n.IP.To4()))
	binary.BigEndian.PutUint32(ip, binary.BigEndian.Uint32(n.IP.To4())|^binary.BigEndian.Uint32(net.IP(n.Mask).To4()))
	return ip, nil
}
