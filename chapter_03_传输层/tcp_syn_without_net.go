package main

// syscall调用RAW socket
// 参考：http://blog.csdn.net/gophers/article/details/20393601
// https://github.com/kdar/gorawtcpsyn

import (
	"bytes"
	"encoding/binary"
	. "fmt"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

type TCPHeader struct {
	SrcPort uint16 // 2^16 最大65535嘛
	DstPort uint16
	SeqNum uint32
	AckNum uint32
	Offset uint8
	Flag uint8
	Window uint16
	Checksum uint16
	UrgentPtr uint16
}

type PsdHeader struct {
	SrcAddr uint32
	DstAddr uint32
	Zero uint8
	ProtoType uint8
	TcpLength uint16
}

// IP 字符串 -> 整数
func inet_addr(ipaddr string) uint32 {
	var (
		segments []string = strings.Split(ipaddr, ".")
		ip [4]uint64 // 为毛是64  - = 而不是8
		ret uint64
	)
	for i:=0; i<4; i++ {
		ip[i], _ = strconv.ParseUint(segments[i], 10, 64)
	}
	ret = ip[3]<<24 + ip[2]<<16 + ip[1]<<8 + ip[0]
	return uint32(ret)
}

// 主机无符号短整形，转化为网络字节顺序。高低位互换
func htons(port uint16) uint16 {
	var (
		high uint16 = port >> 8
		ret uint16 = port << 8 + high
	)
	return ret
}

// 每16bit相加，溢出的回卷，然后取反码
func CheckSum(data []byte) uint16 {
	var (
		sum uint32
		length int = len(data)
		index int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += (sum >> 16) // 回卷
	return uint16(^sum)
}

func main() {
	var (
		msg string
		psdHeader PsdHeader
		tcpHeader TCPHeader
	)

	Printf("Input the content: ")
	Scanf("%s", &msg)

	// 填充TCP伪首部
	psdHeader.SrcAddr = inet_addr("127.0.0.1")
	psdHeader.DstAddr = inet_addr("127.0.0.1")
	psdHeader.Zero = 0
	psdHeader.ProtoType = syscall.IPPROTO_TCP
	psdHeader.TcpLength = uint16(unsafe.Sizeof(TCPHeader{})) + uint16(len(msg))

	// 填充TCP首部
	tcpHeader.SrcPort = htons(3000)
	tcpHeader.DstPort = htons(8080)
	tcpHeader.SeqNum = 0
	tcpHeader.AckNum = 0
	tcpHeader.Offset = uint8(uint16(unsafe.Sizeof(TCPHeader{}))/4) << 4
	tcpHeader.Flag = 2 // SYN
	tcpHeader.Window = 60000
	tcpHeader.Checksum = 0

    // buffer 用来写入两种首部，求得校验和。为什么是两个首部？为什么没有数据？
    var (
    	buffer bytes.Buffer
	)
    binary.Write(&buffer, binary.BigEndian, psdHeader)
    binary.Write(&buffer, binary.BigEndian, tcpHeader)
    tcpHeader.Checksum = CheckSum(buffer.Bytes())

    // 清空buffer，填充要发送的部分
    buffer.Reset()
    binary.Write(&buffer, binary.BigEndian, tcpHeader)
    binary.Write(&buffer, binary.BigEndian, msg)

    // 对RAW socket操作
    var (
    	sockfd int
    	addr syscall.SockaddrInet4
    	err error
	)

    if sockfd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_TCP); err != nil {
    	Println("Socket() error: ", err.Error())
    	return
	}

	defer syscall.Shutdown(sockfd, syscall.SHUT_RDWR)
	addr.Addr[0], addr.Addr[1], addr.Addr[2], addr.Addr[3] = 127, 0, 0, 1
	addr.Port = 8081
	println("Send buffer bytes", buffer.Bytes(), len(buffer.Bytes()))
	if err = syscall.Sendto(sockfd, buffer.Bytes(), 0, &addr); err != nil {
		Println("Sendto() error: ", err.Error())
		return
	}

	Println("Send success!")
}
