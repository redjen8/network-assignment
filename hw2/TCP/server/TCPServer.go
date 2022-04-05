/**
 * TCPServer.go
 **/

package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func main() {
	serverPort := "12000"

	listener, _ := net.Listen("tcp", ":"+serverPort)
	fmt.Printf("Server is ready to receive on port %s\n", serverPort)
	start_time := time.Now()
	request_number := 0
	buffer := make([]byte, 1024)
	command_able := true

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		fmt.Println("Bye bye~")
		os.Exit(0)
	}()

	for {
		conn, _ := listener.Accept()
		fmt.Printf("Connection request from %s\n", conn.RemoteAddr().String())

		for command_able {
			command_able = false
			count, _ := conn.Read(buffer)
			request_number += 1

			recv_message := string(buffer[:count])
			if len(recv_message) == 0 {
				command_able = true
				break
			}
			cmd := strings.ToUpper(recv_message[:11])
			fmt.Println(cmd)
			switch cmd {
			case "ASK_TXTCONV":
				reply_message := strings.ToUpper(recv_message[12:])
				conn.Write([]byte(reply_message))
			case "ASK_IP_PORT":
				reply_message := conn.RemoteAddr().String()
				conn.Write([]byte(reply_message))
			case "ASK_REQ_NUM":
				reply_message := request_number
				conn.Write([]byte(strconv.Itoa(reply_message)))
			case "ASK_RUNTIME":
				reply_message := time.Since(start_time).String()
				conn.Write([]byte(reply_message))
			case "ASK_CONNEND":
				fmt.Println("Bye bye~")
				conn.Close()
				break
			}
			command_able = true
		}
	}
}
