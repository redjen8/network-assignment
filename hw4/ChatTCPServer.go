package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var clientConnInfoMap = map[string]string{}
var clientConnMap = map[string]net.Conn{}

func broadcastMessage(connList map[string]net.Conn, message string) {
	for _, element := range connList {
		element.Write([]byte(message))
	}
}

func connection_handle(conn net.Conn, nickname string) {
	conn_flag := true
	buffer := make([]byte, 1024)

	for conn_flag {
		count, _ := conn.Read(buffer)
		fmt.Println(buffer[:count])
		user_command := string(buffer[0])
		switch user_command {
		case "0":
			fmt.Printf("Welcome User : %s!\n", buffer[1:count])
		case "1":
			// \list command
			fmt.Printf("Welcome User : %s!\n", buffer[1:count])
		case "2":
			// \dm command
			fmt.Printf("")
		case "3":
			// \exit command
			replyMessage := "Client " + nickname + " disconnected.\n"
			fmt.Printf(replyMessage)
			conn_flag = false
			broadcastMessage(clientConnMap, replyMessage)
			delete(clientConnInfoMap, nickname)
			delete(clientConnMap, nickname)
			conn.Close()
		case "4":
			// \ver command
			fmt.Printf("")
		case "5":
			// \rtt command
			fmt.Printf("")
		case "6":
			// without any command, default chat
			chatContent := string(buffer[1:count])
			replyMessage := nickname + " > " + chatContent
			broadcastMessage(clientConnMap, replyMessage)
		}
		fmt.Printf("[DEBUG] command %s : %s from %s\n", user_command, buffer[1:count], nickname+conn.RemoteAddr().String())
	}
}

func main() {

	serverPort := "22848"

	listener, _ := net.Listen("tcp", ":"+serverPort)
	fmt.Printf("Server is ready to receive on port %s\n", serverPort)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		fmt.Println("Bye bye~")
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
		userNickname := string(buffer[1:count])
		fmt.Printf("User nickname : %s from %s connected.\n", userNickname, conn.RemoteAddr())
		if existClientInfo, isAlreadyExists := clientConnInfoMap[userNickname]; isAlreadyExists {
			fmt.Printf("User nickname conflict. User Info : %s\n", existClientInfo)
			_, err = conn.Write([]byte(strconv.Itoa(-1)))
			if err != nil {
				continue
			}
		} else {
			clientConnInfoMap[userNickname] = conn.RemoteAddr().String()
			clientConnMap[userNickname] = conn
			go connection_handle(conn, userNickname)
		}
	}
}
