# import threading 
import socket
import select
import struct
from typing import IO

HEADER = 4
FORMAT = "utf-8"


def create_server(port: int) -> None:
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)

    HOST_IP = socket.gethostbyname(socket.gethostname())
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


def get_content_len(client: socket) -> int:
    content_len = b''
    while len(content_len) < HEADER:
        chunk = client.recv(HEADER - len(content_len))

        content_len += chunk

    content_len = struct.unpack("!I", content_len)
    print(f"done getting content_len -> {content_len[0]}")
    return content_len[0]


def recv_content(client: socket, content_len: int) -> None:
    content = b''
    while True:
        chunk = client.recv(content_len)
        if not chunk:
            print("no more content to read from client")
            break

        content += chunk

    # so the content is an array of [author_len which is 1byte, author_bytes, content_bytes]
    author_len = content[0]
    print( "\n\n", content.decode("utf-8", errors="ignore"))
    # decodec_author = content[1: 1 + author_len].decode(FORMAT)
    # decodec_content = content[1 + author_len:].decode(FORMAT, errors="ignore")

    # print("\nAuthor:", decodec_author)
    # print("Quote: ", decodec_content)
    # print(type(decodec_content))

    # file = open(decodec_author + "co", "w", encoding="utf-8")
    # file.write(decodec_content)
    # file.close()


def handle_conn(client: socket) -> None:
    content_len = get_content_len(client)
    recv_content(client, content_len)


create_server(6000)
