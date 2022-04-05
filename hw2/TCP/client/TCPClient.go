/**
 * TCPClient.go
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

	conn, _ := net.Dial("tcp", serverName+":"+serverPort)

	localAddr := conn.LocalAddr().(*net.TCPAddr)
	fmt.Printf("Client is running on port %d\n", localAddr.Port)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		var message_cmd string = "ASK_CONNEND"
		conn.Write([]byte(message_cmd))
		conn.Close()
		fmt.Println("Bye bye~")
		os.Exit(0)
	}()

	for {
		var user_option int
		fmt.Println("\n<Menu>")
		fmt.Println("Please input your client option in integer")
		fmt.Println("Option 1) Convert text to upper case letters")
		fmt.Println("Option 2) Ask the server what is my ip and port number")
		fmt.Println("Option 3) Ask the server how many requests it has served so far")
		fmt.Println("Option 4) Ask the server how long it has been running since it started")
		fmt.Print("Option :: ")
		fmt.Scanln(&user_option)

		switch user_option {
		case 1:
			var message_cmd, message_text string = "ASK_TXTCONV", ""
			fmt.Print("Input lowercase sentence: ")
			fmt.Scanln(&message_text)
			conn.Write([]byte(message_cmd + ", " + message_text))
			start_time := time.Now()
			buffer := make([]byte, 1024)
			conn.Read(buffer)
			elapsed := time.Since(start_time)
			fmt.Printf("Reply from server: %s\n", string(buffer))
			fmt.Printf("RTT = %d ms\n", elapsed.Milliseconds())
		case 2:
			var message_cmd string = "ASK_IP_PORT"
			conn.Write([]byte(message_cmd))
			start_time := time.Now()
			buffer := make([]byte, 1024)
			conn.Read(buffer)
			elapsed := time.Since(start_time)
			punct_loc := strings.Index(string(buffer), ",")
			fmt.Printf("Reply from server: client IP = %s, port = %s\n", string(buffer)[:punct_loc], string(buffer[punct_loc+1:]))
			fmt.Printf("RTT = %d\n", elapsed.Milliseconds())
		case 3:
			var message_cmd string = "ASK_REQ_NUM"
			conn.Write([]byte(message_cmd))
			start_time := time.Now()
			buffer := make([]byte, 1024)
			conn.Read(buffer)
			elapsed := time.Since(start_time)
			fmt.Printf("Reply from server: requests served = %s\n", string(buffer))
			fmt.Printf("RTT = %d ms\n", elapsed.Milliseconds())
		case 4:
			var message_cmd string = "ASK_RUNTIME"
			conn.Write([]byte(message_cmd))
			start_time := time.Now()
			buffer := make([]byte, 1024)
			conn.Read(buffer)
			elapsed := time.Since(start_time)
			fmt.Printf("Reply from server: run time = %s\n", string(buffer))
			fmt.Printf("RTT = %d ms\n", elapsed.Milliseconds())
		case 5:
			var message_cmd string = "ASK_CONNEND"
			conn.Write([]byte(message_cmd))
			fmt.Println("Bye bye~")
			conn.Close()
			os.Exit(0)
		}
	}
}
