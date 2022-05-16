/**
 * ChatTCPServer.go
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
)

var clientConnInfoMap = map[string]string{}
var clientConnMap = map[string]net.Conn{}
var version = "b0.0.1"

func broadcastMessage(connList map[string]net.Conn, message string, nickname string) {
	if nickname == "" {
		for _, element := range connList {
			element.Write([]byte(message))
		}
	} else {
		for key, element := range connList {
			if key != nickname {
				element.Write([]byte(message))
			}
		}
	}
}

func connection_handle(conn net.Conn, nickname string) {
	conn_flag := true
	buffer := make([]byte, 1024)

	for conn_flag {
		count, _ := conn.Read(buffer)
		user_command := string(buffer[0])
		if count == 0 {
			conn_flag = false
			break
		}
		switch user_command {
		case "1":
			// \list command
			var listCmdReply string
			numClient := 1
			for key, value := range clientConnInfoMap {
				eachClientInfo := key + " " + value
				listCmdReply += "(" + strconv.Itoa(numClient) + ") " + eachClientInfo + "\n"
				numClient += 1
			}
			conn.Write([]byte(listCmdReply))
		case "2":
			// \dm command
			chatContent := string(buffer[1:count])
			lbraceLocation := strings.Index(chatContent, "{")
			rbraceLocation := strings.Index(chatContent, "}")
			if lbraceLocation == -1 || rbraceLocation == -1 {
				conn.Write([]byte("[Error : direct message format is wrong.]"))
			}
			partnerNickname := chatContent[lbraceLocation+1 : rbraceLocation]
			if existClientInfo, isAlreadyExists := clientConnMap[partnerNickname]; isAlreadyExists {
				directMessageContent := "from : " + nickname + "> " + chatContent[rbraceLocation+1:]
				existClientInfo.Write([]byte(directMessageContent))
			} else {
				conn.Write([]byte("[Error : cannot find user '" + partnerNickname + "' on the server.]"))
			}
			if strings.Contains(strings.ToLower(chatContent), "i hate professor") {
				conn.Write([]byte("{err}[You are kicked from chat room.]"))
				delete(clientConnInfoMap, nickname)
				delete(clientConnMap, nickname)
				banMessage := "[" + nickname + " left. There are " + strconv.Itoa(len(clientConnMap)) + " users connected.]"
				broadcastMessage(clientConnMap, banMessage, "")
			}
		case "3":
			// \exit command
			delete(clientConnInfoMap, nickname)
			delete(clientConnMap, nickname)
			replyMessage := "[Client " + nickname + " left. There are " + strconv.Itoa(len(clientConnMap)) + " users connected.]"
			fmt.Println(replyMessage)
			conn_flag = false
			broadcastMessage(clientConnMap, replyMessage, nickname)
		case "4":
			// \ver command
			conn.Write([]byte(version))
		case "5":
			// \rtt command
			rttMessage := "{rtt}" + string(buffer[1:count])
			conn.Write([]byte(rttMessage))
		case "6":
			// without any command, default chat
			chatContent := string(buffer[1:count])
			replyMessage := nickname + " > " + chatContent
			broadcastMessage(clientConnMap, replyMessage, nickname)

			if strings.Contains(strings.ToLower(chatContent), "i hate professor") {
				conn.Write([]byte("{err}[You are kicked from chat room.]"))
				delete(clientConnInfoMap, nickname)
				delete(clientConnMap, nickname)
				banMessage := "[" + nickname + " left. There are " + strconv.Itoa(len(clientConnMap)) + " users connected.]"
				broadcastMessage(clientConnMap, banMessage, "")
			}
		default:
			continue
		}
	}
}

func main() {

	serverPort := "22848"

	listener, _ := net.Listen("tcp", ":"+serverPort)
	fmt.Printf("Server is ready to receive on port %s\n", serverPort)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		fmt.Println("gg~")
		os.Exit(0)
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		defer conn.Close()
		buffer := make([]byte, 1024)
		count, _ := conn.Read(buffer)
		userNickname := string(buffer[1:count])
		if existClientInfo, isAlreadyExists := clientConnInfoMap[userNickname]; isAlreadyExists {
			fmt.Printf("User nickname conflict. User Info : %s\n", existClientInfo)
			_, err = conn.Write([]byte("{err}[that nickname is already used by another user. cannot connect]"))
			if err != nil {
				continue
			}
		} else if len(clientConnMap) >= 8 {
			fmt.Println("User has been blocked, due to max connection limit.")
			_, err = conn.Write([]byte("{err}[chatting room full. cannot connect]"))
			if err != nil {
				continue
			}
		} else {
			clientConnInfoMap[userNickname] = conn.RemoteAddr().String()
			clientConnMap[userNickname] = conn
			go connection_handle(conn, userNickname)
			welcomeMessage := "[" + userNickname + " joined from " + conn.RemoteAddr().String() + ". There are " + strconv.Itoa(len(clientConnMap)) + " users connected.]"
			broadcastMessage(clientConnMap, welcomeMessage, "")
			fmt.Println(welcomeMessage)
		}
	}
}
