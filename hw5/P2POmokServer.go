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

var clientConnInfoMap = map[string]string{}
var clientUDPPortMap = map[string]int{}
var clientConnMap = map[string]net.Conn{}

func connectionHandle(userNickname string) {
	userUDPPort := clientUDPPortMap[userNickname]
	if len(clientConnMap) == 1 {
		// first player connected to server. waits..
		fmt.Println(userNickname + " joined from " + clientConnInfoMap[userNickname] + ". UDP port " + strconv.Itoa(userUDPPort))
		fmt.Println("1 user connected, Waiting for another")
	} else if len(clientConnMap) == 2 {
		// second player connected to server, starts game
		strconv.Itoa(clientUDPPortMap[userNickname])
		fmt.Println(userNickname + " joined from " + clientConnInfoMap[userNickname] + ". UDP port " + strconv.Itoa(userUDPPort))

		var opponentNickname string
		for key, value := range clientConnInfoMap {
			if strings.Compare(userNickname, key) != 0 {
				opponentNickname = value
			}
		}

		fmt.Println("2 user connected, notifying " + opponentNickname + " and " + userNickname)
		playerConnection := clientConnMap[userNickname]
		opponentConnection := clientConnMap[opponentNickname]
		opponentMessage := ""
		opponentConnection.Write([]byte(opponentMessage))

		playerMessage := ""
		playerConnection.Write([]byte(playerMessage))

		delete(clientConnMap, userNickname)
		delete(clientConnMap, opponentNickname)
		delete(clientUDPPortMap, userNickname)
		delete(clientUDPPortMap, opponentNickname)
		delete(clientConnInfoMap, userNickname)
		delete(clientConnInfoMap, opponentNickname)
		fmt.Println(opponentNickname + " and " + userNickname + "disconnected.")
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
		fmt.Println("gg~")
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
		splitIndex := strings.Index(readStr, ".")
		userNickname := readStr[:splitIndex]
		userUDPPort := readStr[splitIndex+1:]

		clientConnInfoMap[userNickname] = conn.RemoteAddr().String()
		clientUDPPortMap[userNickname], _ = strconv.Atoi(userUDPPort)
		clientConnMap[userNickname] = conn
		go connectionHandle(userNickname)
	}
}
