package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var clientConnInfoMap = map[string]string{}
var clientConnMap = map[string]net.Conn{}
var version = "b0.0.1"

func broadcastMessage(connList map[string]net.Conn, message string, nickname string) {
	for key, element := range connList {
		if key != nickname {
			element.Write([]byte(message))
		}
	}
}

func connection_handle(conn net.Conn, nickname string) {
	conn_flag := true
	buffer := make([]byte, 1024)

	for conn_flag {
		count, _ := conn.Read(buffer)
		user_command := string(buffer[0])
		switch user_command {
		case "0":
			fmt.Printf("Welcome User : %s!\n", buffer[1:count])
		case "1":
			// \list command
			for key, value := range clientConnInfoMap {
				eachClientInfo := key + value
				conn.Write([]byte(eachClientInfo))
			}
		case "2":
			// \dm command
			fmt.Printf("")
		case "3":
			// \exit command
			replyMessage := "Client " + nickname + " disconnected.\n"
			fmt.Printf(replyMessage)
			conn_flag = false
			broadcastMessage(clientConnMap, replyMessage, nickname)
			delete(clientConnInfoMap, nickname)
			delete(clientConnMap, nickname)
			conn.Close()
		case "4":
			// \ver command
			conn.Write([]byte(version))
		case "5":
			// \rtt command
			fmt.Printf("")
		case "6":
			// without any command, default chat
			chatContent := string(buffer[1:count])
			replyMessage := nickname + " > " + chatContent
			broadcastMessage(clientConnMap, replyMessage, nickname)
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
			_, err = conn.Write([]byte("that nickname is already used by another user. cannot connect"))
			if err != nil {
				continue
			}
		} else if len(clientConnMap) >= 8 {
			fmt.Printf("User %s has been blocked, due to max connection limit.\n", existClientInfo)
			_, err = conn.Write([]byte("chatting room full. cannot connect"))
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
