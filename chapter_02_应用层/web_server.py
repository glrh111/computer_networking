#! /usr/bin/env python
# coding: utf-8

"""
使用socket开发一个HTTP服务器

接收GET请求。并返回一个静态资源文件。

GET localhost:8000/path
如果当前目录下path存在，那么返回文件内容；否则返回404
"""


from socket import *
import os
import datetime

sock = socket(AF_INET, SOCK_STREAM)
sock.bind(('127.0.0.1', 8000))
sock.listen(0)
print 'Start server at localhost:8000...\n\n'


def parse_data(data):
    """解析出路径，请求方法等参数
    return (method, path) or (None, None)
    GET /path HTTP/1.1
    Host: localhost
    """
    line_list = data.split('\r\n')
    try:
        method, path, _  = line_list[0].split(' ')
        path = path[1:]
    except Exception:
        method, path = None, None
    return method, path


def find_body_string(path):
    """如果找到这个文件，返回文件内容；否则返回404notfound
    return: (is_found, return_string)
    """

    if path and os.path.isfile(path):
        with open(path, 'r') as f:
            return_string = f.read()
            is_found = True
    else:
        return_string = '404 Not Found'
        is_found = False

    return is_found, return_string


# 判断接收的内容
try:
    while True:
        conn, (addr, port) = sock.accept()
        print '[{}]: Connect to client from {}:{}'.format(datetime.datetime.now().isoformat(), addr, port)

        data = ''

        while '\r\n\r\n' not in data:
            data += conn.recv(1000) # 接收数据的时候需要判断，是否接收完了
            print data

        print '接收完毕：', data

        method, path = parse_data(data)
        print 'method:', method, 'path:', path
        is_found, body_string = find_body_string(path)

        status_code = None
        reason_phase = None

        if is_found:
            status_code = 200
            reason_phase = "OK"
        else:
            status_code = 404
            reason_phase = "Not Found"

        response = 'HTTP/1.1 {status_code} {reason_phase}\r\n\r\n{response_body_string}\r\n'.format(**{
            "status_code": status_code,
            "reason_phase": reason_phase,
            "response_body_string": body_string
        })

        conn.send(response)
        conn.close()
except KeyboardInterrupt:
    print 'Close server!\n'
    sock.close()







