package main

// syscall调用RAW socket

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
	SrcPort uint32 // 为毛是32,32位？TCP中该字段是16位
	DstPort uint32
	Zero uint8
	ProtoType uint8
	


}