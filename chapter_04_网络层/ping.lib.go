package main

// 参考 https://github.com/sparrc/go-ping/blob/master/ping.go
// 参考 https://github.com/tatsushid/go-fastping/blob/master/fastping.go

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"math"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	timeSliceLength  = 8
	protocolICMP     = 1
	protocolIPv6ICMP = 58
)

var (
	ipv4Proto = map[string]string{"ip": "ip4:icmp", "udp": "udp4"}
	ipv6Proto = map[string]string{"ip": "ip6:ipv6-icmp", "udp": "udp6"}
)

func NewPinger(addr string) (*Pinger, error) {
	ipaddr, err := net.ResolveIPAddr("ip", addr)
	if err != nil {
		return nil, err
	}

	var ipv4I bool
	if isIPv4(ipaddr.IP) {
		ipv4I = true
	} else if isIPv6(ipaddr.IP) {
		ipv4I = false
	}

	return &Pinger{
		Interval:    time.Second,
		Timeout:     time.Second * 100000, // 这个地方，如果用 10^6会有问题？什么原因，我曹
		Count:       -1,
		Debug:       false,
		PacketsSent: 0,
		PacketsRecv: 0,
		done:        make(chan bool),
		ipaddr:      ipaddr,
		addr:        addr,
		ipv4:        ipv4I,
		network:     "ip",
		size:        timeSliceLength,
	}, nil
}

// pinger
type Pinger struct {
	// Interval is the wait time between each packet send. Default is 1s
	Interval time.Duration

	// Timeout specifies a timeout before ping exits, regardless of how many packets have been received
	Timeout time.Duration

	// Count tells pinger to stop after sending (and receiving) Count echo
	// packets. If this option is not specified, pinger will operate until interrupted.
	Count int

	// Debug runs in debug mode
	Debug bool

	// Number of packets sent
	PacketsSent int

	// Number of packets received
	PacketsRecv int

	// rtts is all of the rtts
	rtts []time.Duration

	// OnRecv is called when pinger receives and processed a packet
	OnRecv func(*Packet)

	// OnFinish is called when pinger exits
	OnFinish func(*Statistics)

	// stop chan bool
	done chan bool

	//
	ipaddr *net.IPAddr
	addr   string

	ipv4     bool
	source   string
	size     int
	sequence int
	network  string
}

type packet struct {
	bytes  []byte
	nbytes int
}

// Packet represents a received and processed ICMP echo packet.
type Packet struct {
	// rtt is the round-trip time it took to ping
	Rtt time.Duration

	// IPAddr is the address of host being pinged
	IPAddr *net.IPAddr

	// Nbytes is the number os bytes in the message
	Nbytes int

	// Seq is the ICMP sequence number
	Seq int
}

//
type Statistics struct {
	PacketsRecv int

	PacketsSent int

	PacketsLoss float64

	IPAddr *net.IPAddr

	Addr string

	Rtts []time.Duration

	TotalRtt time.Duration

	MinRtt time.Duration

	MaxRtt time.Duration

	AvgRtt time.Duration

	// 标准差
	StdDevRtt time.Duration
}

// SetIPAddr sets the ip address of the target host
func (p *Pinger) SetIPAddr(ipaddr *net.IPAddr) {
	var ifIPv4 bool
	if isIPv4(ipaddr.IP) {
		ifIPv4 = true
	} else if isIPv6(ipaddr.IP) {
		ifIPv4 = false
	}

	p.ipaddr = ipaddr
	p.addr = ipaddr.String()
	p.ipv4 = ifIPv4
}

// IPAddr return the ip address of the target host
func (p *Pinger) IPAddr() *net.IPAddr {
	return p.ipaddr
}

// SetAddr resolves and sets the ip address of the target host
func (p *Pinger) SetAddr(addr string) error {
	ipaddr, err := net.ResolveIPAddr("ip", addr)
	if err != nil {
		return err
	}
	p.SetIPAddr(ipaddr)
	p.addr = addr
	return nil
}

// Addr return the string ip address of the target host
func (p *Pinger) Addr() string {
	return p.addr
}

// SetPrivileged sets the type of ping pinger will send
func (p *Pinger) SetPrivileged(privileged bool) {
	if privileged {
		p.network = "ip"
	} else {
		p.network = "udp"
	}
}

// Privileged return the property
func (p *Pinger) Privileged() bool {
	return "ip" == p.network
}

// run
func (p *Pinger) run() {
	var conn *icmp.PacketConn
	if p.ipv4 {
		// network 可以是 ip, udp; ipv4Proto 返回 ip4:icmp, udp4; ip6:ipv6-icmp, udp6
		if conn = p.listen(ipv4Proto[p.network], p.source); conn == nil {
			return
		}
	} else {
		if conn = p.listen(ipv6Proto[p.network], p.source); conn == nil {
			return
		}
	}

	defer conn.Close()
	defer p.finish()

	var wg sync.WaitGroup
	recv := make(chan *packet, 5)
	wg.Add(1)
	go p.recvICMP(conn, recv, &wg)

	err := p.sendICMP(conn)
	if err != nil {
		fmt.Println(err.Error())
	}

	timeout := time.NewTicker(p.Timeout)
	interval := time.NewTicker(p.Interval)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	for {
		select {
		case <-c:
			close(p.done)
		case <-p.done:
			wg.Wait()
			return
		case <-timeout.C:
			close(p.done)
			wg.Wait()
			return
		case <-interval.C:
			err = p.sendICMP(conn)
			if err != nil {
				fmt.Println("FATAL: ", err.Error())
			}
		case r := <-recv:
			err := p.processPacket(r)
			if err != nil {
				fmt.Println("FATAL: ", err.Error())
			}
		default:
			if p.Count > 0 && p.PacketsRecv >= p.Count {
				close(p.done)
				wg.Wait()
				return
			}

		}
	}

}

// Run
func (p *Pinger) Run() {
	p.run()
}

func (p *Pinger) finish() {
	handler := p.OnFinish
	if handler != nil {
		s := p.Statistics()
		handler(s)
	}
}

func (p *Pinger) Statistics() *Statistics {
	loss := float64(p.PacketsSent-p.PacketsRecv) / float64(p.PacketsSent) * 100
	var min, max, total time.Duration
	if len(p.rtts) > 0 {
		min = p.rtts[0]
		max = p.rtts[0]
	}
	for _, rtt := range p.rtts {
		if min > rtt {
			min = rtt
		}
		if max < rtt {
			max = rtt
		}
		total += rtt
	}
	s := Statistics{
		PacketsRecv: p.PacketsRecv,
		PacketsSent: p.PacketsSent,
		PacketsLoss: loss,
		IPAddr:      p.ipaddr,
		Addr:        p.addr,
		Rtts:        p.rtts,
		MinRtt:      min,
		MaxRtt:      max,
		TotalRtt:    total,
	}
	if len(p.rtts) > 0 {
		s.AvgRtt = total / time.Duration(len(p.rtts))
		stdDevRtt := time.Duration(0)
		for _, rtt := range p.rtts {
			stdDevRtt += (rtt - s.AvgRtt) * (rtt - s.AvgRtt)
		}
		s.StdDevRtt = time.Duration(math.Sqrt(
			float64(stdDevRtt / time.Duration(len(s.Rtts)))))
	}
	return &s

}

func (p *Pinger) recvICMP(
	conn *icmp.PacketConn,
	recv chan<- *packet,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	for {
		select {
		case <-p.done:
			return
		default:
			bytes := make([]byte, 512)
			conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
			n, _, err := conn.ReadFrom(bytes)
			if err != nil {
				if neterr, ok := err.(*net.OpError); ok {
					if neterr.Timeout() {
						// Read timeout
						continue
					} else {
						close(p.done)
						return
					}
				}
			}
			recv <- &packet{bytes: bytes, nbytes: n}
		}
	}
}

func (p *Pinger) processPacket(recv *packet) error {
	var bytes []byte
	var proto int
	if p.ipv4 {
		if p.network == "ip" {
			bytes = ipv4Payload(recv.bytes)
		} else {
			bytes = recv.bytes
		}
		proto = protocolICMP
	} else {
		bytes = recv.bytes
		proto = protocolIPv6ICMP
	}
	var m *icmp.Message
	var err error
	if m, err = icmp.ParseMessage(proto, bytes[:recv.nbytes]); err != nil {
		return fmt.Errorf("Error parsing ICMP message")
	}

	if m.Type != ipv4.ICMPTypeEchoReply && m.Type != ipv6.ICMPTypeEchoReply {
		// 不是echo reply
		return nil
	}

	outPkt := &Packet{
		Nbytes: recv.nbytes,
		IPAddr: p.ipaddr,
	}

	switch pkt := m.Body.(type) {
	case *icmp.Echo:
		outPkt.Rtt = time.Since(bytesToTime(pkt.Data[:timeSliceLength]))
		outPkt.Seq = pkt.Seq
		p.PacketsRecv += 1
	default:
		// vary bad, not sure how it could happen
		return fmt.Errorf("Error, invalid ICMP echo reply. Body bytes: %T, %s", pkt, pkt)
	}
	p.rtts = append(p.rtts, outPkt.Rtt)
	handler := p.OnRecv
	if handler != nil {
		handler(outPkt)
	}
	return nil

}

func (p *Pinger) sendICMP(conn *icmp.PacketConn) error {
	var typ icmp.Type
	if p.ipv4 {
		typ = ipv4.ICMPTypeEcho
	} else {
		typ = ipv6.ICMPTypeEchoRequest
	}

	var dst net.Addr = p.ipaddr
	if p.network == "udp" {
		dst = &net.UDPAddr{IP: p.ipaddr.IP, Zone: p.ipaddr.Zone}
	}

	t := timeToBytes(time.Now())
	if p.size-timeSliceLength != 0 {
		t = append(t, byteSliceOfSize(p.size-timeSliceLength)...)
	}
	bytes, err := (&icmp.Message{
		Type: typ, Code: 0,
		Body: &icmp.Echo{
			ID:   rand.Intn(65535),
			Seq:  p.sequence,
			Data: t,
		},
	}).Marshal(nil)

	if err != nil {
		return err
	}

	for {
		if _, err := conn.WriteTo(bytes, dst); err != nil {
			if neterr, ok := err.(*net.OpError); ok {
				if neterr.Err == syscall.ENOBUFS {
					continue
				}
			}
		}
		p.PacketsSent += 1
		p.sequence += 1
		break
	}

	return nil

}

// listen
func (p *Pinger) listen(netProto string, source string) *icmp.PacketConn {
	conn, err := icmp.ListenPacket(netProto, source)
	if err != nil {
		fmt.Printf("Error listening for ICMP packets: %s\n", err.Error())
		close(p.done)
		return nil
	}
	return conn
}

func byteSliceOfSize(n int) []byte {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = 1
	}
	return b

}

// b[0]: 版本(4bit)+首部长度(4bit)+服务类型 一共16bit
// b[0]的后4bit，表示首部长度，单位为32位字长
func ipv4Payload(b []byte) []byte {
	if len(b) < ipv4.HeaderLen {
		return b
	}
	hdrlen := int(b[0]&0x0f) << 2
	return b[hdrlen:]
}

// 这个函数搞毛？
func bytesToTime(b []byte) time.Time {
	var nsec int64
	for i := uint8(0); i < 8; i++ {
		nsec += int64(b[i]) << ((7 - i) * 8)
	}
	return time.Unix(nsec/1000000000, nsec%1000000000)
}

func timeToBytes(t time.Time) []byte {
	nsec := t.UnixNano()
	b := make([]byte, 8)
	for i := uint8(0); i < 8; i++ {
		b[i] = byte((nsec >> ((7 - i) * 8)) & 0xff)
	}
	return b
}

func isIPv4(ip net.IP) bool {
	return len(ip.To4()) == net.IPv4len
}

func isIPv6(ip net.IP) bool {
	return len(ip) == net.IPv6len
}
