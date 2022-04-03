#
# SimpleEchoTCPServer.py
#

from socket import *

serverPort = 12000
serverSocket = socket(AF_INET, SOCK_STREAM)
serverSocket.bind(('', serverPort))
serverSocket.listen(1)
requestNumber = 0

print("Server is ready to receive on port", serverPort)

while True:
    (connectionSocket, clientAddress) = serverSocket.accept()
    print('Connection request from', clientAddress)
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
            replyMessage = 'runtime'
        elif cmd == 'ASK_CONNEND':
            print('Bye bye~')
            connectionSocket.close()
        else:
            replyMessage = 'Wrong Command'

        connectionSocket.send(replyMessage.encode())

    except KeyboardInterrupt:
        print('Bye bye~')
        break

    connectionSocket.close()

