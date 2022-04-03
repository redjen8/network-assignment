#
# SimpleEchoTCPServer.py
#

from socket import *
import time

serverPort = 12000
serverSocket = socket(AF_INET, SOCK_STREAM)
serverSocket.bind(('', serverPort))
serverSocket.listen(1)
requestNumber = 0

print("Server is ready to receive on port", serverPort)

while True:
    runTime = time.time()
    (connectionSocket, clientAddress) = serverSocket.accept()
    print('Connection request from', clientAddress)
    command_able = True
    while command_able:
        command_able = False
        try:
            message = connectionSocket.recv(2048).decode()
            requestNumber += 1

            cmd = message.upper()[:11]
            replyMessage = ''

            if cmd == 'ASK_TXTCONV':
                replyMessage = message[12:].upper()
            elif cmd == 'ASK_IP_PORT':
                replyMessage = clientAddress[0] + ',' + str(clientAddress[1])
            elif cmd == 'ASK_REQ_NUM':
                replyMessage = str(requestNumber)
            elif cmd == 'ASK_RUNTIME':
                currentTime = time.time()
                replyMessage = str(currentTime - runTime)
            elif cmd == 'ASK_CONNEND':
                print('Bye bye~')
                command_able = True
                break
                
            else:
                replyMessage = 'Wrong Command'
                command_able = True
                break

            connectionSocket.send(replyMessage.encode())

        except KeyboardInterrupt:
            print('Bye bye~')
            break
        except ConnectionAbortedError:
            print('Bye bye~')
            break

        command_able = True

    connectionSocket.close()

