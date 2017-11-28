#! /usr/bin/env python
# coding: utf-8

"""
写一个pingpong程序，客户端发送ping，服务端返回pong，并记录时间
如果长时间未响应，那么返回miss
"""

from socket import *
import datetime
import time
import random


def log(message):
    print '[{}] {}'.format(datetime.datetime.now().isoformat(), message)


class UDPPingPongServer(object):

    def deal_with_message(self):
        data, (addr, port) = self.sock.recvfrom(1000)
        log('Received message from {}:{}: [{}]'.format(addr, port, data))
        # time.sleep(random.random())
        self.sock.sendto("pong", (addr, port))

    def run_server(self):
        self.sock = socket(AF_INET, SOCK_DGRAM)
        self.sock.bind(('127.0.0.1', 8000))
        log('Start server at 127.0.0.1:8000\n')

        try:
            while True:
                self.deal_with_message()
        except KeyboardInterrupt:
            self.sock.close()
            log('Closed server!')


if __name__ == '__main__':
    UDPPingPongServer().run_server()


