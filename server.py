import socket
import threading 
"""
Threading module is used to handle concurrency in python. In socket programming 
there are tons of blocking code that can inhibit "async" operations, if thats your speed.
Python itself is single threaded, ie it can only run one task at a time, but the threading 
module in python solves that by using actual OS-Level threads.

Some blocking tasks in socket module include:
    - listening for connections ie server.listen() and server.accept()
    - reading from client ie server.recv() and server.sendmsg()
"""


HEADER = 64
FORMAT = "utf-8"
DISCONNECT_MSG = "!DISCONNECTED"
port = 5050
# automatically get IP_ADDRESS of the running device
SERVER = socket.gethostbyname(socket.gethostname())
ADDR = (SERVER, port)

server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server.bind(ADDR)


def handle_client(conn: socket, addr: str):
    print(f"[new_client] -> connected from {addr}")

    connected = True
    while connected:
        msg_len = conn.recv(HEADER).decode(FORMAT)
        
        if len(msg_len) < 1:
            print("client did not send anything")
            break

        msg_len = int(msg_len)
        msg = conn.recv(msg_len).decode(FORMAT)

        if msg == "!DISCONNECTED":
            connected = False
        print(f"{addr} ->> {msg}")

    conn.close()


def start():
    server.listen()
    while True:
        conn, addr = server.accept()
        thread = threading.Thread(target=handle_client, args=(conn, addr))
        thread.start()

        print(f"[active_conns] -> {threading.active_count() - 1}")


print(f"[starting] -> server started at {SERVER}")
start()
