/**
 * TCPClient.go
 **/

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
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
		fmt.Print("Option :: ")
		fmt.Scanln(&user_option)

		switch user_option {
		case 1:
			var message_cmd, message_text string = "ASK_TXTCONV", ""
			fmt.Print("Input lowercase sentence: ")
			conn.Write([]byte(message_cmd + ", " + message_text))
			start_time := time.Now()
			buffer := make([]byte, 1024)
			conn.Read(buffer)
			elapsed_time := time.Since(start_time)
			fmt.Printf("Reply from server: %s", string(buffer))
			fmt.Printf("RTT = %s", string(elapsed_time))
		}
	}

	fmt.Printf("Input lowercase sentence: ")
	input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	conn.Write([]byte(input))

	buffer := make([]byte, 1024)
	conn.Read(buffer)
	fmt.Printf("Reply from server: %s", string(buffer))

	conn.Close()
}
