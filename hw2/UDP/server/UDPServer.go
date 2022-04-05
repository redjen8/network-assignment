/**
 * UDPServer.go
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

	pconn, _ := net.ListenPacket("udp", ":"+serverPort)
	fmt.Printf("Server is ready to receive on port %s\n", serverPort)
	start_time := time.Now()
	request_number := 0
	buffer := make([]byte, 1024)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		fmt.Println("Bye bye~")
		os.Exit(0)
	}()

	for {
		count, r_addr, _ := pconn.ReadFrom(buffer)
		fmt.Printf("Connection requested from %s %d\n", r_addr, count)
		request_number += 1
		recv_message := string(buffer[:count])
		cmd := strings.ToUpper(recv_message[:11])
		switch cmd {
		case "ASK_TXTCONV":
			reply_message := strings.ToUpper(recv_message[12:])
			fmt.Println(reply_message)
			pconn.WriteTo([]byte(reply_message), r_addr)
		case "ASK_IP_PORT":
			reply_message := r_addr.String()
			fmt.Println(reply_message)
			pconn.WriteTo([]byte(reply_message), r_addr)
		case "ASK_REQ_NUM":
			reply_message := request_number
			pconn.WriteTo([]byte(strconv.Itoa(reply_message)), r_addr)
		case "ASK_RUNTIME":
			reply_message := time.Since(start_time).String()
			pconn.WriteTo([]byte(reply_message), r_addr)
		case "ASK_CONNEND":
			fmt.Println("Bye bye~")
		}
	}
}
