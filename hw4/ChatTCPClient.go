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

func sendTCPData(command_int int, conn net.Conn) {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(command_int))
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

	localAddr := conn.LocalAddr().(*net.TCPAddr)
	fmt.Printf("Welcome %s to CAU network classroom at %s:%s. Client is running on port %d.\n", nickname, serverName, serverPort, localAddr.Port)

	go func() {
		<-signals
		fmt.Println("Bye bye~")
		os.Exit(0)
	}()

	go readServerUpdate(conn)

	for {
		var user_input string
		fmt.Scanln(&user_input)
		slice := strings.Split(user_input, " ")

		switch slice[0] {
		case "\\list":
			command_int := 1
			sendTCPData(command_int, conn)
		case "\\dm":
			command_int := 2
			sendTCPData(command_int, conn)
		case "\\exit":
			command_int := 3
			sendTCPData(command_int, conn)
		case "\\ver":
			command_int := 4
			sendTCPData(command_int, conn)
		case "\\rtt":
			command_int := 5
			sendTCPData(command_int, conn)
		default:
			command_int := 0
			sendTCPData(command_int, conn)
		}
		break
	}
}
