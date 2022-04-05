/**
 * TCPClient.go
 **/

package main

import (
	"fmt"
	"net"
	"time"
)

func main() {

	serverName := "localhost"
	serverPort := "12000"

	conn, _ := net.Dial("tcp", serverName+":"+serverPort)

	localAddr := conn.LocalAddr().(*net.TCPAddr)
	fmt.Printf("Client is running on port %d\n", localAddr.Port)

	for {
		var user_option int
		fmt.Println("<Menu>")
		fmt.Println("Please input your client option in integer")
		fmt.Println("Option 1) Convert text to upper case letters")
		fmt.Println("Option 2) Ask the server what is my ip and port number")
		fmt.Println("Option 3) Ask the server how many requests it has served so far")
		fmt.Println("Option 4) Ask the server how long it has been running since it started")
		fmt.Println("Option :: ")
		fmt.Scanln(&user_option)

		switch user_option {
		case 1:
			var message_cmd, message_text string = "ASK_TXTCONV", ""
			fmt.Print("Input lowercase sentence: ")
			conn.Write([]byte(message_cmd + ", " + message_text))
			start_time := time.Now()
			buffer := make([]byte, 1024)
			conn.Read(buffer)
			elapsed := time.Since(start_time)
			fmt.Printf("Reply from server: %s", string(buffer))
			fmt.Printf("RTT = %s\n", elapsed)
		case 2:
			var message_cmd string = "ASK_IP_PORT"
			conn.Write([]byte(message_cmd))
			start_time := time.Now()
			buffer := make([]byte, 1024)
			conn.Read(buffer)
			elapsed := time.Since(start_time)
			fmt.Printf("Reply from server: client IP = %s, port = %s", string(buffer), string(buffer))
			fmt.Printf("RTT = %s ms\n", elapsed)
		}
	}

	conn.Close()
}
