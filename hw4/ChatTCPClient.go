package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

func sendTCPData(command_int int, message string, conn net.Conn) {
	if len(message) == 0 {
		message = strconv.Itoa(command_int)
	} else {
		message = strconv.Itoa(command_int) + message
	}
	conn.Write([]byte(message))
}

func readServerUpdate(conn net.Conn) {
	for {
		buffer := make([]byte, 1024)
		count, _ := conn.Read(buffer)
		response := string(buffer[:count])
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
	if len(os.Args) >= 2 {
		nickname = os.Args[1]
	} else {
		nickname = "redjen"
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	conn, err := net.Dial("tcp", serverName+":"+serverPort)
	sendTCPData(0, nickname, conn)

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
		user_input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		if strings.Contains(user_input, "\r\n") {
			user_input = user_input[:len(user_input)-2]
		} else if strings.Contains(user_input, "\n") {
			user_input = user_input[:len(user_input)-1]
		}
		fmt.Println()
		input_slice := strings.Split(user_input, " ")
		commandRegex, _ := regexp.Compile("\\\\(\\w)+")

		if commandRegex.MatchString(input_slice[0]) {
			command := input_slice[0]
			//need to fix- what if len(command) < 3 and dm ?
			//if len(command) < 3 {
			//	fmt.Println("Invalid Command")
			//}
			command = command[1:len(command)]
			switch command {
			case "list":
				sendTCPData(1, "", conn)
			case "dm":
				partnerNickname := input_slice[1]
				partnerMessage := ""
				for idx, eachMessage := range input_slice {
					if idx <= 1 {
						continue
					}
					partnerMessage += eachMessage + " "
				}
				sendTCPData(2, "{"+partnerNickname+"}"+partnerMessage, conn)
			case "exit":
				sendTCPData(3, "", conn)
				clientFlag = false
			case "ver":
				sendTCPData(4, "", conn)
			case "rtt":
				sendTCPData(5, "", conn)
			default:
				fmt.Println("Invalid Command")
			}
		} else {
			sendTCPData(6, user_input, conn)
		}
	}
	fmt.Println("gg~")
}
