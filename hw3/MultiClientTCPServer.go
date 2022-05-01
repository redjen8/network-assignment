/**
 * EasyTCPServer.go
 * 20172848 Jeong Seok Woo
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

func connection_handle(conn net.Conn, name int, conn_list map[int]bool, start_time time.Time, request_number int) {
	conn_flag := true
	buffer := make([]byte, 1024)

	for conn_flag {
		count, _ := conn.Read(buffer)
		request_number += 1

		recv_message := string(buffer[:count])
		if len(recv_message) == 0 {
			fmt.Println("Command Length 0, Invalid!")
		}
		cmd := strings.ToUpper(recv_message[:11])

		switch cmd {
		case "ASK_TXTCONV":
			reply_message := strings.ToUpper(recv_message[12:])
			_, err := conn.Write([]byte(reply_message))
			if err != nil {
				return
			}
			fmt.Printf("Command %s Executed From Client %d(%s).\n", strconv.Itoa(request_number), name, conn.RemoteAddr().String())
		case "ASK_IP_PORT":
			reply_message := conn.RemoteAddr().String()
			_, err := conn.Write([]byte(reply_message))
			if err != nil {
				return
			}
			fmt.Printf("Command %s Executed From Client %d(%s).\n", strconv.Itoa(request_number), name, conn.RemoteAddr().String())
		case "ASK_REQ_NUM":
			reply_message := request_number
			_, err := conn.Write([]byte(strconv.Itoa(reply_message)))
			if err != nil {
				return
			}
			fmt.Printf("Command %s Executed From Client %d(%s).\n", strconv.Itoa(request_number), name, conn.RemoteAddr().String())
		case "ASK_RUNTIME":
			reply_message := time.Time{}.Add(time.Since(start_time))
			_, err := conn.Write([]byte(reply_message.Format("15:04:05")))
			if err != nil {
				return
			}
			fmt.Printf("Command %s Executed From Client %d(%s).\n", strconv.Itoa(request_number), name, conn.RemoteAddr().String())
		case "ASK_CONNEND":
			fmt.Printf("Client %d disconnected. Number of connected clients = %d\n", name, len(conn_list)-1)
			err := conn.Close()
			if err != nil {
				return
			}
			delete(conn_list, name)
			conn_flag = false
			break
		}
	}
}

func add(s map[int]bool, v int) map[int]bool {
	s[v] = true
	return s
}

func main() {
	serverPort := "22848"

	listener, _ := net.Listen("tcp", ":"+serverPort)
	fmt.Printf("Server is ready to receive on port %s\n", serverPort)
	start_time := time.Now()
	request_number := 0
	connection_number := 0
	connection_list := map[int]bool{}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		fmt.Println("Bye bye~")
		os.Exit(0)
	}()

	ticker := time.NewTicker(60 * time.Second)
	timeout := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Printf("Current Number of Active Clients : %d, Current Time : %s\n", len(connection_list), time.Now())
			case <-timeout:
				ticker.Stop()
				return
			}
		}
	}()

	for {
		conn, _ := listener.Accept()
		connection_number += 1
		connection_list = add(connection_list, connection_number)
		fmt.Printf("Client %d connected. Number of connected clients = %d\n", connection_number, len(connection_list))
		go connection_handle(conn, connection_number, connection_list, start_time, request_number)
	}
}
