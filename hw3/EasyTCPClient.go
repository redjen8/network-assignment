/**
 * EasyTCPClient.go
 * 20172848 Jeong Seok Woo
 **/

package main

import (
	"bufio"
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
	serverPort := "22848"

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
		fmt.Println("Option 5) Exit")
		fmt.Print("Option :: ")
		fmt.Scanln(&user_option)

		switch user_option {
		case 1:
			var message_cmd, message_text string = "ASK_TXTCONV", ""
			fmt.Print("Input lowercase sentence: ")
			message_text, _ = bufio.NewReader(os.Stdin).ReadString('\n')
			buffer := make([]byte, 1024)
			start_time := time.Now()
			conn.Write([]byte(message_cmd + "," + message_text))
			count, _ := conn.Read(buffer)
			elapsed := float64((time.Now().Sub(start_time)).Nanoseconds()) / 1000000
			fmt.Printf("Reply from server: %s\n", string(buffer[:count]))
			fmt.Printf("RTT = %f ms\n", elapsed)
		case 2:
			var message_cmd string = "ASK_IP_PORT"
			buffer := make([]byte, 1024)
			start_time := time.Now()
			conn.Write([]byte(message_cmd))
			count, _ := conn.Read(buffer)
			elapsed := float64((time.Now().Sub(start_time)).Nanoseconds()) / 1000000
			punct_loc := strings.LastIndex(string(buffer), ":")
			fmt.Printf("Reply from server: client IP = %s, port = %s\n", string(buffer)[:punct_loc], string(buffer[punct_loc+1:count]))
			fmt.Printf("RTT = %f ms\n", elapsed)
		case 3:
			var message_cmd string = "ASK_REQ_NUM"
			buffer := make([]byte, 1024)
			start_time := time.Now()
			conn.Write([]byte(message_cmd))
			count, _ := conn.Read(buffer)
			elapsed := float64((time.Now().Sub(start_time)).Nanoseconds()) / 1000000
			fmt.Printf("Reply from server: requests served = %s\n", string(buffer[:count]))
			fmt.Printf("RTT = %f ms\n", elapsed)
		case 4:
			var message_cmd string = "ASK_RUNTIME"
			buffer := make([]byte, 1024)
			start_time := time.Now()
			conn.Write([]byte(message_cmd))
			count, _ := conn.Read(buffer)
			elapsed := float64((time.Now().Sub(start_time)).Nanoseconds()) / 1000000
			fmt.Printf("Reply from server: run time = %s\n", string(buffer[:count]))
			fmt.Printf("RTT = %f ms\n", elapsed)
		case 5:
			var message_cmd string = "ASK_CONNEND"
			conn.Write([]byte(message_cmd))
			fmt.Println("Bye bye~")
			conn.Close()
			os.Exit(0)
		}
	}
}
