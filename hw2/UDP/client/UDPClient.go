/**
 * UDPClient.go
 **/

package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	serverName := "localhost"
	serverPort := "12000"

	pconn, _ := net.ListenPacket("udp", ":")

	localAddr := pconn.LocalAddr().(*net.UDPAddr)
	fmt.Printf("Client is running on port %d\n", localAddr.Port)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		var message_cmd string = "ASK_CONNEND"
		server_addr, _ := net.ResolveUDPAddr("udp", serverName+":"+serverPort)
		pconn.WriteTo([]byte(message_cmd), server_addr)
		pconn.Close()
		fmt.Println("Bye bye~")
		os.Exit(0)
	}()

	for {
		fmt.Println("\n<Menu>")
		fmt.Println("Please input your client option in integer")
		fmt.Println("Option 1) Convert text to upper case letters")
		fmt.Println("Option 2) Ask the server what is my ip and port number")
		fmt.Println("Option 3) Ask the server how many requests it has served so far")
		fmt.Println("Option 4) Ask the server how long it has been running since it started")
		fmt.Print("Option :: ")

		var user_option int
		fmt.Scanln(&user_option)
		buffer := make([]byte, 1024)

		switch user_option {
		case 1:
			var message_cmd, message_text string = "ASK_TXTCONV", ""
			fmt.Print("Input lowercase sentence: ")
			fmt.Scanln(&message_text)
			server_addr, _ := net.ResolveUDPAddr("udp", serverName+":"+serverPort)
			pconn.WriteTo([]byte(message_cmd+","+message_text), server_addr)
			start_time := time.Now()
			count, _, _ := pconn.ReadFrom(buffer)
			elapsed := time.Since(start_time)
			fmt.Printf("Reply from server: %s\n", string(buffer[:count]))
			fmt.Printf("RTT = %d ms\n", elapsed.Milliseconds())
		case 2:
			var message_cmd string = "ASK_IP_PORT"
			server_addr, _ := net.ResolveUDPAddr("udp", serverName+":"+serverPort)
			pconn.WriteTo([]byte(message_cmd), server_addr)
			start_time := time.Now()
			count, _, _ := pconn.ReadFrom(buffer)
			elapsed := time.Since(start_time)
			punct_loc := strings.Index(string(buffer), ":")
			fmt.Printf("Reply from server: client IP = %s, port = %s\n", string(buffer)[:punct_loc], string(buffer[punct_loc+1:count]))
			fmt.Printf("RTT = %d ms\n", elapsed.Milliseconds())
		case 3:
			var message_cmd string = "ASK_REQ_NUM"
			server_addr, _ := net.ResolveUDPAddr("udp", serverName+":"+serverPort)
			pconn.WriteTo([]byte(message_cmd), server_addr)
			start_time := time.Now()
			count, _, _ := pconn.ReadFrom(buffer)
			elapsed := time.Since(start_time)
			fmt.Printf("Reply from server: requests served = %s\n", string(buffer[:count]))
			fmt.Printf("RTT = %d ms\n", elapsed.Milliseconds())
		case 4:
			var message_cmd string = "ASK_RUNTIME"
			server_addr, _ := net.ResolveUDPAddr("udp", serverName+":"+serverPort)
			pconn.WriteTo([]byte(message_cmd), server_addr)
			start_time := time.Now()
			count, _, _ := pconn.ReadFrom(buffer)
			elapsed := time.Since(start_time)
			fmt.Printf("Reply from server: run time = %s\n", string(buffer[:count]))
			fmt.Printf("RTT = %d ms\n", elapsed.Milliseconds())
		case 5:
			var message_cmd string = "ASK_CONNEND"
			server_addr, _ := net.ResolveUDPAddr("udp", serverName+":"+serverPort)
			pconn.WriteTo([]byte(message_cmd), server_addr)
			fmt.Println("Bye bye~")
			pconn.Close()
			os.Exit(0)
		}
	}
}
