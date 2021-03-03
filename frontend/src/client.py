import socket
import errno
import time
import sys
import os
import json

HOST = '127.0.0.1'  # The server's hostname or IP address
PORT = 8080        # The port used by the server
BUFFER_SIZE = 1024

with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:

    s.connect((HOST, PORT))
    s.setblocking(0)

    while(True):
        try:
            print("input command, 0-9")
            a = sys.stdin.readline()
            s.sendall((a + "\n").encode())

            data = s.recv(BUFFER_SIZE)
            data = json.loads(data)
            print("received data: ", data)
        except socket.error as e:
            if e.args[0] == errno.EWOULDBLOCK:
                time.sleep(1)
            else:
                print(e)
                break
#        print("received data: ", data["tanks"]["0"]["x"])

    s.close()
