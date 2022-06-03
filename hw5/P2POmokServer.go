package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var clientConnInfoMap = map[string]string{}
var clientConnMap = map[string]net.Conn{}

func connectionHandle(userNickname string) {

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

		clientConnInfoMap[userNickname] = conn.RemoteAddr().String() + ":" + userUDPPort
		clientConnMap[userNickname] = conn
		go connectionHandle(userNickname)
	}
}
