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

func connection_handle(conn net.Conn) {
	conn_flag := true
	buffer := make([]byte, 1024)

	for conn_flag {
		count, _ := conn.Read(buffer)
		user_command := buffer[0]
		switch user_command {
		case 0:
			fmt.Printf("Welcome User : %s!\n", buffer[1])
		case 1:
			// \list command
			fmt.Printf("Welcome User : %s!\n", buffer[1])
		case 2:
			// \dm command
			fmt.Printf("")
		case 3:
			// \exit command
			fmt.Printf("")
			conn_flag = false
		case 4:
			// \ver command
			fmt.Printf("")
		case 5:
			// \rtt command
			fmt.Printf("")
		case 99:
			// without any command, default chat
			fmt.Printf("")
		}
		fmt.Printf("command %d : %s from %s\n", user_command, buffer[:count], conn.RemoteAddr().String())
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
		userNickname := string(buffer[:count])
		fmt.Printf("User nickname : %s\n", userNickname)
		fmt.Printf("User addr : %s\n", conn.RemoteAddr())
		if existClientInfo, isAlreadyExists := clientConnInfoMap[userNickname]; isAlreadyExists {
			fmt.Printf("User nickname conflict. User Info : %s\n", existClientInfo[0])
			_, err = conn.Write([]byte(strconv.Itoa(-1)))
			if err != nil {
				continue
			}
		} else {
			clientConnInfoMap[userNickname] = conn.RemoteAddr().String()
			clientConnMap[userNickname] = conn
			go connection_handle(conn)
		}
	}
}
