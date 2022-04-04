#
# SimpleEchoUDPServer.py
#

from socket import *
import datetime
import time

serverPort = 12000
serverSocket = socket(AF_INET, SOCK_DGRAM)
serverSocket.bind(('', serverPort))
requestNumber = 0
byeMessage = 'Bye bye~'

print("Server is ready to receive on port", serverPort)
runTime = time.time()

while True:
    try:
        message, clientAddress = serverSocket.recvfrom(2048)
        print('Connection requested from', clientAddress)

        message = message.decode()
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
            replyMessage = str(datetime.timedelta(seconds=currentTime-runTime)).split(".")[0]
        elif cmd == 'ASK_CONNEND':
            print(byeMessage)
            
        else:
            replyMessage = 'Wrong Command'

        serverSocket.sendto(replyMessage.encode(), clientAddress)

    except KeyboardInterrupt:
        print(byeMessage)
        break
    except ConnectionAbortedError:
        print(byeMessage)
    
    
