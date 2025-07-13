import socket 
import struct 
# import os 
from typing import IO

HOST = socket.gethostbyname(socket.gethostname())
PORT = 6000
ADDR = (HOST, PORT)
FORMAT = "utf-8"


def create_client() -> socket:
    try:
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        s.connect(ADDR)
        return s

    except socket.error as e:
        print(f"An error occured: {e}")
        return None


def file_handler(fname: str, flag="r") -> IO:
    try:
        file = open(fname, flag)
        return file
    except Exception as e:
        print(f"An eror occured: {e}")
        exit()


def ack_header(client: socket, n: int) -> None:
    # we first going to create a header that tells the server:
    # HEY! Im going to send you N_PACKETS, here's the author/name of the file
    # If you say YES, I UNDERSTAND, we are going to start a countdown, just to know that 
    # we are on the same page
    # For this protocol we might use delimiters, tired of the infinite slicing 
    """
    ====================================
         CLIENT: Acknowledge Request
    ====================================
    Type: Acknowledge \r\n
    n_packets: 59 \r\n
    start_at: 0 \r\n
    Author: server.go \r\n

    ====================================
         SERVER: Acknowledge Response
    ====================================
    Type: Acknowledge \r\n
    Status: Acknowledged \r\n


    ====================================
         CLIENT: Process
    ====================================
    Type: Data-Packet \r\n
    sync: 1 \r\n
    curr-at: 1 \r\n
    ack: 0 \r\n
    sent: 1\r\n
    \r\n
    b'data content'

    Since we already told the server to be expecting 59 packets, we can add headers along lines of:
    I acknowledge the 0th data-packet you sent, 
    I've synced my count by 1, and 
    I'm waiting_on the next 1 and I've got 1 packet(s)
    ====================================
         SERVER: Process
    ====================================
    acked: 0 \r\n
    sync: 1 \r\n
    wait_on: 1 \r\n
    packets-recived: 1 \r\n
    """

    pass


def stream_file(fname: str) -> None:
    # s = create_client()
    file = file_handler(fname)

    content = file.readlines()

    codec_author = fname.encode(FORMAT)
    codec_content = [line.encode(FORMAT) for line in content]
    N_PACKETS = len(codec_content)

    codec_author_len = struct.pack("B", len(codec_author))

    base_body = codec_author_len + codec_author 
    base_body_len = len(base_body)

    packet_sync = 0

    while packet_sync < N_PACKETS:
        header = struct.pack("!I", (len(codec_content[packet_sync]))) 
        data_packet = header + codec_content[packet_sync]

        # s.sendall(packet)
        print(f"new packet -> {packet_sync}\n packet-details: {data_packet}")
        packet_sync += 1

    print(f"file_contents:\n{codec_content}")


stream_file("server.py")
