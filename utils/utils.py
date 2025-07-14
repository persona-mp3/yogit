import socket
from typing import IO


def create_client(ADDR: tuple) -> socket:
    try:
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)

        s.connect(ADDR)
        return s
    except Exception as e:
        print(f"An error occured in creating server: \n {e}")
        exit()


def file_handler(fname: str, flag="r") -> IO:
    try:
        f = open(fname, flag)
        return f
    except Exception as e:
        print(f"An error occured in file_handler\n{e}")
        exit()
    pass
