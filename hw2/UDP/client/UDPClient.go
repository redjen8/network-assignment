/**
 * UDPClient.go
 **/

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	serverName := "localhost"
	serverPort := "12000"

	pconn, _ := net.ListenPacket("udp", ":")

	localAddr := pconn.LocalAddr().(*net.UDPAddr)
	fmt.Printf("Client is running on port %d\n", localAddr.Port)

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

		switch user_option {
		case 1:
			var message_cmd, message_text string = "ASK_TXTCONV", ""
			fmt.Print("Input lowercase sentence: ")
			fmt.Scanln(&message_text)
			server_addr, _ := net.ResolveUDPAddr("udp", serverName+":"+serverPort)
			pconn.WriteTo([]byte(message_cmd+", "+message_text), server_addr)
			start_time := time.Now()
			buffer := make([]byte, 1024)
			pconn.ReadFrom(buffer)
			elapsed := time.Since(start_time)
			fmt.Printf("Reply from server: %s\n", string(buffer))
			fmt.Printf("RTT = %d ms\n", elapsed.Milliseconds())
			pconn.Close()
		case 2:
			var message_cmd string = "ASK_IP_PORT"
			pconn.WriteTo([]byte(message_cmd))
			start_time := time.Now()
			buffer := make([]byte, 1024)
			pconn.ReadFrom(buffer)
			elapsed := time.Since(start_time)
			punct_loc := strings.Index(string(buffer), ",")
			fmt.Printf("Reply from server: client IP = %s, port = %s\n", string(buffer)[:punct_loc], string(buffer[punct_loc+1:]))
			fmt.Printf("RTT = %d\n", elapsed.Milliseconds())
		case 3:
		case 4:
		case 5:

		}
	}

	fmt.Printf("Input lowercase sentence: ")
	input, _ := bufio.NewReader(os.Stdin).ReadString('\n')

	server_addr, _ := net.ResolveUDPAddr("udp", serverName+":"+serverPort)
	pconn.WriteTo([]byte(input), server_addr)

	buffer := make([]byte, 1024)
	pconn.ReadFrom(buffer)
	fmt.Printf("Reply from server: %s", string(buffer))
}
