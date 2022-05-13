package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func sendTCPData(command_int int, message string, conn net.Conn) {
	buf := make([]byte, 1024)
	binary.LittleEndian.PutUint32(buf, uint32(command_int))
	if message != "" {
		copy(buf[1:], message)
		fmt.Printf(string(buf[:]))
	} else {
		fmt.Printf(string(buf[:]))
	}
	conn.Write(buf)
}

func readServerUpdate(conn net.Conn) {
	for {
		buffer := make([]byte, 1024)
		count, _ := conn.Read(buffer)
		fmt.Printf("Reply from server: %s\n", string(buffer[:count]))
	}
}

func main() {

	serverName := "localhost"
	serverPort := "22848"

	nickname := os.Args[1]

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	conn, _ := net.Dial("tcp", serverName+":"+serverPort)
	sendTCPData(0, nickname, conn)

	conn.Write([]byte(nickname))

	localAddr := conn.LocalAddr().(*net.TCPAddr)
	fmt.Printf("Welcome %s to CAU network classroom at %s:%s. Client is running on port %d.\n", nickname, serverName, serverPort, localAddr.Port)

	go func() {
		<-signals
		fmt.Println("Bye bye~")
		os.Exit(0)
	}()

	go readServerUpdate(conn)

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
		sendTCPData(99, "some message", conn)
	}
}
