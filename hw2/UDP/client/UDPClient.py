#
# SimpleEchoUDPClient.py
#

from socket import AF_INET, SOCK_DGRAM, socket
import time

serverName = 'localhost'
serverPort = 12000

clientSocket = socket(AF_INET, SOCK_DGRAM)
clientSocket.bind(('', 0))

print("Client is running on port", clientSocket.getsockname()[1])

while True:
    print('\n<Menu>')
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
            clientSocket.sendto((message_cmd + ',' + message_text).encode(), (serverName, serverPort))
            start_time = time.perf_counter()
            modifiedMessage, serverAddress = clientSocket.recvfrom(2048)
            elapsed_time = time.perf_counter() - start_time
            print('Reply from server:', modifiedMessage.decode())
            print(f'RTT = {elapsed_time*1000:.2f} ms')

        elif option == 2:
            # send 'ASK_IP_PORT'
            # recv [xxx.xxx.xxx.xxx, XXXX] which is my ip and port
            message = 'ASK_IP_PORT'
            clientSocket.sendto(message.encode(), (serverName, serverPort))
            start_time = time.perf_counter()
            modifiedMessage, serverAddress = clientSocket.recvfrom(2048)
            elapsed_time = time.perf_counter()-start_time
            modifiedMessage = modifiedMessage.decode()
            punct_loc = modifiedMessage.find(',')
            print(f'Reply from server: client IP = {modifiedMessage[:punct_loc]}, port = {modifiedMessage[punct_loc+1:]}')
            print(f'RTT = {elapsed_time*1000:.2f} ms')

        elif option == 3:
            # send 'ASK_REQ_NUM'
            # recv [xxx] which is count of requests served
            message = 'ASK_REQ_NUM'
            clientSocket.sendto(message.encode(), (serverName, serverPort))
            start_time = time.perf_counter()
            modifiedMessage, serverAddress = clientSocket.recvfrom(2048)
            elapsed_time = time.perf_counter()-start_time
            print('Reply from server: requests served = ', modifiedMessage.decode())
            print(f'RTT = {elapsed_time*1000:.2f} ms')

        elif option == 4:
            # send 'ASK_RUNTIME'
            # recv [xxx] which is run time of the server since it has been started
            message = 'ASK_RUNTIME'
            clientSocket.sendto(message.encode(), (serverName, serverPort))
            start_time = time.perf_counter()
            modifiedMessage, serverAddress = clientSocket.recvfrom(2048)
            elapsed_time = time.perf_counter()-start_time
            print('Reply from server: run time = ', modifiedMessage.decode())
            print(f'RTT = {elapsed_time*1000:.2f} ms')

        elif option == 5:
            message = 'ASK_CONNEND'
            clientSocket.sendto(message.encode(), (serverName, serverPort))
            print('Bye bye~')
            break

        else:
            print('Wrong option Selected. Please Select menu again.')

    except KeyboardInterrupt:
        message = 'ASK_CONNEND'
        clientSocket.sendto(message.encode(), (serverName, serverPort))
        print('Bye bye~')
        break
    
    except ValueError:
        print('Wrong Input, Please Input Again.')

clientSocket.close()
