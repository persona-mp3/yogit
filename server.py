# import threading 
import socket
import select
import struct
from typing import Sequence
from utils.utils import file_handler

HEADER = 4
FORMAT = "utf-8"
ACK = "Acknowledge"

TYPE_DATA_TRANS = "Packet"


def create_server(port: int) -> None:
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)

    # HOST_IP = socket.gethostbyname(socket.gethostname())
    HOST_IP = "127.0.0.1"
    ADDR = (HOST_IP, port)

    s.bind(ADDR)
    s.listen(5)

    print(f"server crearted, acitive @ {ADDR}")
    while True:
        read_conns, _, _ = select.select([s], [], [], 1)
        if read_conns:
            for conns in read_conns:
                conn, addr = conns.accept()
                print(f"[new-client] -> @ {addr}")
                try:
                    handle_conn(conn)
                except socket.error as e:
                    print(f"An error occured in socket operation: {e}")
                    continue
                finally:
                    conn.close()


#  the function could return a string that describes the kind of request made by a client
def recv_ack_header(c: socket) -> Sequence[str]:
    content_len = b''
    while len(content_len) < HEADER:
        chunk = c.recv(HEADER - len(content_len))
        content_len += chunk

    content_len = struct.unpack("!I", content_len)[0]
    print(f"content length of ack request: {content_len}")

    content = b''
    while len(content) < content_len:
        chunk = c.recv(content_len - len(content))
        content += chunk

    print("\n\ndecoding header...\n\n") 

    decoded_content = content.decode(FORMAT)
    req_body = decoded_content.split()

    if req_body[0] != ACK:        # return (req_body[0], req_body)
        return 

    print(f"Request-Type: {req_body[0]}")
    print(f"Packets-to-recieve: {req_body[1]}")
    print(f"Author: {req_body[2]}")
    print(f"Client-start: {req_body[3]}")

    Type = "Acknowledge \r\n"
    Status = "Acknowledged \r\n"

    encoded_response = Type.encode(FORMAT) + Status.encode(FORMAT)
    c.sendall(encoded_response)
    print()
    n_packets = req_body[1]
    author = req_body[2]
    return (n_packets, author)


def recv_content(c: socket, author: str) -> None:
    content_len = b''
    while len(content_len) < HEADER:
        chunk = c.recv(HEADER - len(content_len))
        content_len += chunk

    content_len = struct.unpack("!I", content_len)[0]
    # print(f"content-len: {content_len}")

    content = b''
    while len(content) < content_len:
        chunk = c.recv(content_len - len(content))
        content += chunk

    # print(content)
    type_req_len_bytes = content[:1]
    type_req_len = struct.unpack("B", type_req_len_bytes)[0]

    type_req = content[1:type_req_len + 1]

    if type_req.decode(FORMAT) != TYPE_DATA_TRANS:
        print(f"this packet is not a data transmission type, {type_req}")
        return

    # print(f"len of this header-type: {type_req_len}, type_request: {type_req}")

    data = content[type_req_len + 1:]
    data = data.decode(FORMAT)
    print(data)

    f = file_handler(("./test/" + author), "a")
    f.write(data)


def handle_conn(client: socket) -> None:
    details = recv_ack_header(client)
    n_packets = int(details[0])
    author = str(details[1])

    packet_sync = 0
    while packet_sync < n_packets:
        if packet_sync == n_packets:
            print("done over here, not allowed to send more than proposed")
            break
        recv_content(client, author)
        packet_sync += 1
        # print(f"pack-sync -> {packet_sync}")


create_server(6000)
