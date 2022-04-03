#
# SimpleEchoTCPClient.py
#

from socket import *

serverName = 'localhost'
serverPort = 12000

clientSocket = socket(AF_INET, SOCK_STREAM)
clientSocket.connect((serverName, serverPort))

print("Client is running on port", clientSocket.getsockname()[1])

while True:
    print('<Menu>')
    print('Please input your client option in integer')
    print('Option 1) Convert text to upper case letters')
    print('Option 2) Ask the server what is my ip and port number')
    print('Option 3) Ask the server how many requests it has served so far')
    print('Option 4) Ask the server how long it has been running since it started')
    try:
        option = int(input('Option :: '))

        if option == 1:
            # send 'ASK_TXTCONV, msg'
            # recv ['converted_msg']
            message_cmd = 'ASK_TXTCONV'
            message_text = input('Input lowercase sentence: ')
            clientSocket.send((message_cmd + ',' + message_text).encode())
            modifiedMessage = clientSocket.recv(2048)
            print('Reply from server:', modifiedMessage.decode())

        elif option == 2:
            # send 'ASK_IP_PORT'
            # recv [xxx.xxx.xxx.xxx, XXXX] which is my ip and port
            message = 'ASK_IP_PORT'
            clientSocket.send(message.encode())
            modifiedMessage = clientSocket.recv(2048).decode()
            print('Reply from server: client IP = {}, port = {}', )

        elif option == 3:
            # send 'ASK_REQ_NUM'
            # recv [xxx] which is count of requests served
            message = 'ASK_REQ_NUM'
            clientSocket.send(message.encode())
            modifiedMessage = clientSocket.recv(2048)
            print('Reply from server: requests served = ', modifiedMessage.decode())

        elif option == 4:
            # send 'ASK_RUNTIME'
            # recv [xxx] which is run time of the server since it has been started
            message = 'ASK_RUNTIME'
            clientSocket.send(message.encode())
            modifiedMessage = clientSocket.recv(2048)
            print('Reply from server: run time = ', modifiedMessage.decode())

        elif option == 5:
            print('Bye bye~')

        else:
            print('Wrong Option Selected. Please Select menu again.')

    except KeyboardInterrupt:
        print('Bye bye~')
        break

clientSocket.close()
