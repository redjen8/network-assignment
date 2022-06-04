package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

var clientConnInfoMap = map[string]string{} // stores ip address of client
var clientUDPPortMap = map[string]int{}     // stores opened udp port of client
var clientConnMap = map[string]net.Conn{}   // stores actual tcp connection of client

func connectionHandle(userNickname string) {
	userUDPPort := clientUDPPortMap[userNickname]
	var opponentNickname string

	if len(clientConnMap) == 1 {
		// first player connected to server. waits..
		fmt.Println(userNickname + " joined from " + clientConnInfoMap[userNickname] + ". UDP port " + strconv.Itoa(userUDPPort) + ".")
		fmt.Println("1 user connected, Waiting for another")
	} else if len(clientConnMap) == 2 {
		// second player connected to server, starts game
		strconv.Itoa(clientUDPPortMap[userNickname])
		fmt.Println(userNickname + " joined from " + clientConnInfoMap[userNickname] + ". UDP port " + strconv.Itoa(userUDPPort) + ".")

		for key, _ := range clientConnInfoMap {
			if strings.Compare(userNickname, key) != 0 {
				opponentNickname = key
			}
		}

		fmt.Println("2 user connected, notifying " + opponentNickname + " and " + userNickname)
		playerConnection := clientConnMap[userNickname]
		opponentConnection := clientConnMap[opponentNickname]
		opponentMessage := "0" + userNickname + "," + clientConnInfoMap[userNickname] + ":" + strconv.Itoa(clientUDPPortMap[userNickname])
		opponentConnection.Write([]byte(opponentMessage))

		// first player came into server plays first
		playerMessage := "1" + opponentNickname + "," + clientConnInfoMap[opponentNickname] + ":" + strconv.Itoa(clientUDPPortMap[opponentNickname])
		playerConnection.Write([]byte(playerMessage))

		delete(clientConnMap, userNickname)
		delete(clientUDPPortMap, userNickname)
		delete(clientConnInfoMap, userNickname)

		delete(clientConnMap, opponentNickname)
		delete(clientUDPPortMap, opponentNickname)
		delete(clientConnInfoMap, opponentNickname)
		fmt.Println(opponentNickname + " and " + userNickname + " disconnected.")
	}
}

func main() {
	serverPort := "52848"

	listener, _ := net.Listen("tcp", ":"+serverPort)
	fmt.Printf("Server is ready to receive on port %s\n", serverPort)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		fmt.Println("bye bye~")
		os.Exit(0)
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		defer conn.Close()

		buffer := make([]byte, 1024)
		count, _ := conn.Read(buffer)

		readStr := string(buffer[:count])
		splitIndex := strings.Index(readStr, ":")
		userNickname := readStr[:splitIndex]
		userUDPPort := readStr[splitIndex+1:]

		if existClientInfo, isAlreadyExists := clientConnInfoMap[userNickname]; isAlreadyExists {
			fmt.Printf("User nickname conflict. User Info : %s\n", existClientInfo)
			_, err = conn.Write([]byte("that nickname is already used by another user. cannot connect"))
			if err != nil {
				continue
			}
		}
		ipAddrIdx := strings.LastIndex(conn.RemoteAddr().String(), ":")
		clientConnInfoMap[userNickname] = conn.RemoteAddr().String()[:ipAddrIdx]
		clientUDPPortMap[userNickname], _ = strconv.Atoi(userUDPPort)
		clientConnMap[userNickname] = conn
		go connectionHandle(userNickname)
	}
}
