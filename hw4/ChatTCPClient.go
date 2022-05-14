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

func sendTCPData(command_int int, message string, conn net.Conn) {
	if len(message) == 0 {
		message = strconv.Itoa(command_int)
		fmt.Printf(message)
	} else {
		message = strconv.Itoa(command_int) + message
		fmt.Printf(message)
	}
	conn.Write([]byte(message))
}

func readServerUpdate(conn net.Conn) {
	for {
		buffer := make([]byte, 1024)
		count, _ := conn.Read(buffer)
		response := string(buffer[:count])
		fmt.Println(response)
		fmt.Printf("Reply from server: %s\n", response)
	}
}

func detectTCPDisconnect(signals chan os.Signal, err error) {
	for {
		if err != nil {
			fmt.Println("TCP Connection Disconnected.")
			signal.Notify(signals, syscall.SIGTERM)
		}
	}
}

func main() {

	serverName := "localhost"
	serverPort := "22848"

	nickname := ""
	if len(os.Args) > 2 {
		nickname = os.Args[1]
	} else {
		nickname = "redjen"
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	conn, err := net.Dial("tcp", serverName+":"+serverPort)
	sendTCPData(0, nickname, conn)

	conn.Write([]byte(nickname))

	localAddr := conn.LocalAddr().(*net.TCPAddr)
	fmt.Printf("Welcome %s to CAU network classroom at %s:%s. Client is running on port %d.\n", nickname, serverName, serverPort, localAddr.Port)

	go func() {
		<-signals
		sendTCPData(3, "", conn)
		fmt.Println("gg~")
		os.Exit(0)
	}()

	go readServerUpdate(conn)
	go detectTCPDisconnect(signals, err)

	clientFlag := true
	for clientFlag {
		var user_input string
		fmt.Scanln(&user_input)
		slice := strings.Split(user_input, " ")

		switch slice[0] {
		case "\\list":
			sendTCPData(1, "", conn)
		case "\\dm":
			sendTCPData(2, "message", conn)
		case "\\exit":
			sendTCPData(3, "", conn)
			clientFlag = false
		case "\\ver":
			sendTCPData(4, "", conn)
		case "\\rtt":
			sendTCPData(5, "", conn)
		}
		sendTCPData(6, slice[0], conn)
	}
}
