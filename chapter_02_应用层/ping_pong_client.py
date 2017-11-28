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

def time13():
    return int(1000*time.time())


class UDPPingPongClient(object):

    def __init__(self):
        self.rtt_list = []

    def send_message(self):
        self.sock = socket(AF_INET, SOCK_DGRAM)
        start_time13 = time13()
        self.sock.sendto("ping", ('127.0.0.1', 8000))
        data, (addr, port) = self.sock.recvfrom(1000)
        end_time13 = time13()
        rtt = end_time13-start_time13
        self.rtt_list.append(rtt)
        log("Received [{}] from {}:{} RTT {}ms".format(data, addr, port, rtt))
        self.sock.close()

    def pingpingping(self):
        log('Start ping to 127.0.0.1:8000\n')

        try:
            while True:
                self.send_message()
                time.sleep(0.5)
        except KeyboardInterrupt:
            log('\n\nPing total amount: {}\nPing average RTT: {}ms\n'.format(
                len(self.rtt_list), float(sum(self.rtt_list))/len(self.rtt_list)
            ))
            log('Stop ping!')


if __name__ == '__main__':
    UDPPingPongClient().pingpingping()


