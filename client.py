import socket 
import struct 

HOST = socket.gethostbyname(socket.gethostname())
PORT = 59999
ADDR = (HOST, PORT)
FORMAT = "utf-8"


def test_receive_conn() -> None:
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    s.connect(ADDR)

    message = "Commes Des Garcons CDG Sil Vu plait - Issac Newton"
    author = "ISSAC NEWTON"

    encoded_msg = message.encode(FORMAT)
    encoded_auth = author.encode(FORMAT)

    author_len = len(encoded_auth)

    packed_author = struct.pack("B", author_len)
    body = packed_author + encoded_auth + encoded_msg 

    header = struct.pack("!I", len(body))

    # so the final_packet is an array of 
    # [content_len(4), author_len(1), author_bytes, content_bytes]
    packet = header + body

    s.sendall(packet)
    s.close()

    print("connection closed")


test_receive_conn()
