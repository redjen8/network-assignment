package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	Row = 10
	Col = 10
)

type Board [][]int

var isPlayerTurn, isPlayerBlack, isGameOnProgress bool
var board Board
var win int
var timer = time.NewTicker(time.Second * 10)

func printBoard(b Board) {
	fmt.Print("   ")
	for j := 0; j < Col; j++ {
		fmt.Printf("%2d", j)
	}

	fmt.Println()
	fmt.Print("  ")
	for j := 0; j < 2*Col+3; j++ {
		fmt.Print("-")
	}

	fmt.Println()

	for i := 0; i < Row; i++ {
		fmt.Printf("%d |", i)
		for j := 0; j < Col; j++ {
			c := b[i][j]
			if c == 0 {
				fmt.Print(" +")
			} else if c == 1 {
				fmt.Print(" @")
			} else if c == 2 {
				fmt.Print(" 0")
			} else {
				fmt.Print(" |")
			}
		}

		fmt.Println(" |")
	}

	fmt.Print("  ")
	for j := 0; j < 2*Col+3; j++ {
		fmt.Print("-")
	}

	fmt.Println()
}

func checkWin(b Board, x, y int) int {
	lastStone := b[x][y]
	startX, startY, endX, endY := x, y, x, y

	// Check X
	for startX-1 >= 0 && b[startX-1][y] == lastStone {
		startX--
	}
	for endX+1 < Row && b[endX+1][y] == lastStone {
		endX++
	}

	if endX-startX+1 >= 5 {
		return lastStone
	}

	// Check Y
	startX, startY, endX, endY = x, y, x, y
	for startY-1 >= 0 && b[x][startY-1] == lastStone {
		startY--
	}
	for endY+1 < Row && b[x][endY+1] == lastStone {
		endY++
	}

	if endY-startY+1 >= 5 {
		return lastStone
	}

	// Check Diag 1
	startX, startY, endX, endY = x, y, x, y
	for startX-1 >= 0 && startY-1 >= 0 && b[startX-1][startY-1] == lastStone {
		startX--
		startY--
	}
	for endX+1 < Row && endY+1 < Col && b[endX+1][endY+1] == lastStone {
		endX++
		endY++
	}

	if endY-startY+1 >= 5 {
		return lastStone
	}

	// Check Diag 2
	startX, startY, endX, endY = x, y, x, y
	for startX-1 >= 0 && endY+1 < Col && b[startX-1][endY+1] == lastStone {
		startX--
		endY++
	}
	for endX+1 < Row && startY-1 >= 0 && b[endX+1][startY-1] == lastStone {
		endX++
		startY--
	}

	if endY-startY+1 >= 5 {
		isGameOnProgress = false
		return lastStone
	}

	return 0
}

func clear() {
	//fmt.Printf("%s", runtime.GOOS)

	clearMap := make(map[string]func()) //Initialize it
	clearMap["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clearMap["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	value, ok := clearMap[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                             //if we defined a clearMap func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clearMap terminal screen :(")
	}
}

func handleUDPConnection(udpConn net.PacketConn, opponentNickname string) {
	buffer := make([]byte, 1024)
	for {
		count, _, _ := udpConn.ReadFrom(buffer)
		readMessages := string(buffer[:count])
		messageCommand := readMessages[0:1]
		switch messageCommand {
		case "5":
			// normal chat message
			fmt.Println(opponentNickname + "> " + readMessages[1:])
		case "2":
			// omok game command
			commandMessage := strings.Split(readMessages, " ")
			newPosX := commandMessage[1]
			newPosY := commandMessage[2]
			opponentColor, _ := strconv.Atoi(commandMessage[3])
			x, _ := strconv.Atoi(newPosX)
			y, _ := strconv.Atoi(newPosY)
			board[x][y] = opponentColor
			clear()
			printBoard(board)
			win = checkWin(board, x, y)
			if win != 0 {
				fmt.Printf("you lose\n")
				isGameOnProgress = false
				break
			}
			count += 1
			if count == Row*Col {
				fmt.Printf("draw!\n")
				isGameOnProgress = false
				break
			}
			isPlayerTurn = true
			if isGameOnProgress {
				timer = time.NewTicker(10 * time.Second)
				go omokPlayTimer(udpConn)
			}
		case "3":
			// gg command
			if isGameOnProgress {
				fmt.Println("you win")
				fmt.Println("Opponent called gg.")
				isGameOnProgress = false
			}
		case "4":
			// exit command
			if isGameOnProgress {
				fmt.Println("you win")
				isGameOnProgress = false
			}
			fmt.Println("Opponent Exited.")
		}

	}
}

func omokPlayTimer(udpConn net.PacketConn) {
	for {
		select {
		case <-timer.C:
			fmt.Println("10 seconds elapsed!!!!")
			timer.Stop()
			return
		}
	}
}

func main() {

	userNickname := ""
	if len(os.Args) >= 2 {
		userNickname = os.Args[1]
	} else {
		userNickname = "redjen"
	}

	serverName := "localhost"
	serverPort := "52848"

	tcpConn, _ := net.Dial("tcp", serverName+":"+serverPort)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		tcpConn.Close()
		fmt.Println("Bye~")
		os.Exit(0)
	}()

	fmt.Println("Welcome " + userNickname + " to p2p-omok server at " + serverName + ":" + serverPort + ".")
	udpConn, _ := net.ListenPacket("udp", ":")
	udpConnPort := udpConn.LocalAddr().(*net.UDPAddr).Port
	welcomeMessage := userNickname + ":" + strconv.Itoa(udpConnPort)
	tcpConn.Write([]byte(welcomeMessage))

	buffer := make([]byte, 1024)
	tcpCount, _ := tcpConn.Read(buffer)
	readFromBuffer := string(buffer[1:tcpCount])
	opponentEndpointIdx := strings.LastIndex(readFromBuffer, ".")
	opponentNickname := readFromBuffer[:opponentEndpointIdx]
	opponentEndpoint, _ := net.ResolveUDPAddr("udp", readFromBuffer[opponentEndpointIdx+1:])

	if strings.Compare(string(buffer[0]), "0") == 0 {
		isPlayerTurn = true
		isPlayerBlack = true
		fmt.Println(opponentNickname + " joined " + opponentEndpoint.String() + ", you play first.\n")
	} else {
		fmt.Println(opponentNickname + " is waiting for you (" + opponentEndpoint.String() + ").")
		isPlayerTurn = false
		isPlayerBlack = false
		fmt.Println(opponentNickname + " plays first.\n")
	}
	tcpConn.Close()
	isGameOnProgress = true
	var colorStr string
	if isPlayerBlack {
		colorStr = "1"
	} else {
		colorStr = "2"
	}

	go handleUDPConnection(udpConn, opponentNickname)

	x, y, count, win := -1, -1, 0, 0
	for i := 0; i < Row; i++ {
		var tempRow []int
		for j := 0; j < Col; j++ {
			tempRow = append(tempRow, 0)
		}
		board = append(board, tempRow)
	}
	printBoard(board)

	for {
		userInput, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		if strings.Contains(userInput, "\r\n") {
			userInput = strings.TrimRight(userInput, "\r\n")
		} else if strings.Contains(userInput, "\n") {
			userInput = strings.TrimSuffix(userInput, "\n")
		}

		// if player's turn, get standard input from client
		commandRegex, _ := regexp.Compile("\\\\(\\w)+")
		eachWord := strings.Split(userInput, " ")
		if eachWord[0] == "\\\\" {
			if isPlayerTurn {
				// coordinates, move on to omok game
				if len(eachWord) != 3 {
					fmt.Println("error, must enter x y!")
					time.Sleep(1 * time.Second)
					continue
				}
				x, _ = strconv.Atoi(eachWord[1])
				y, _ = strconv.Atoi(eachWord[2])
				if x < 0 || y < 0 || x >= Row || y >= Col {
					fmt.Println("error, out of bound!")
					time.Sleep(1 * time.Second)
					continue
				} else if board[x][y] != 0 {
					fmt.Println("error, already used!")
					time.Sleep(1 * time.Second)
					continue
				}

				if isPlayerBlack {
					board[x][y] = 1
				} else {
					board[x][y] = 2
				}
				clear()
				printBoard(board)
				win = checkWin(board, x, y)
				if win != 0 {
					fmt.Printf("you win\n")
					isGameOnProgress = false
				}
				count += 1
				if count == Row*Col {
					fmt.Printf("draw!\n")
					isGameOnProgress = false
				}
				isPlayerTurn = false
				sendMessage := "2 " + eachWord[1] + " " + eachWord[2] + " " + colorStr
				// 1 byte to identify omok command, using "2"
				udpConn.WriteTo([]byte(sendMessage), opponentEndpoint)
				timer.Stop()
			} else {
				fmt.Println("It's not your turn!")
			}
		} else if commandRegex.MatchString(eachWord[0]) {
			// other commands
			// \gg and \exit branch
			command := eachWord[0][1:]
			switch command {
			case "gg":
				// 1 byte to identify gg command, using "3"
				if isGameOnProgress {
					sendMessage := "3"
					udpConn.WriteTo([]byte(sendMessage), opponentEndpoint)
					isGameOnProgress = false
				} else {
					fmt.Println("Game already finished. Ignore command..")
				}
			case "exit":
				// 1 byte to identify exit command, using "4"
				fmt.Println("Bye~")
				sendMessage := "4"
				udpConn.WriteTo([]byte(sendMessage), opponentEndpoint)
				os.Exit(0)
			default:
				fmt.Println("Invalid Command")
			}
		} else {
			udpConn.WriteTo([]byte("5"+userInput), opponentEndpoint)
		}
	}
}
