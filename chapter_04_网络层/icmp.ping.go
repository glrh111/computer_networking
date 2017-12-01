package main

// 已知BUG1 network设置为 udp，用不了
// TTL 需要从ipheader里边读取，暂时没有这个功能

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("Error: no hostname to ping.\n")
		return
	}

	pinger, err := NewPinger(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
	pinger.OnRecv = func(packet *Packet) {
		fmt.Printf("%d bytes from %s (%s): icmp_seq=%d ttl=[None] time=%v \n",
			packet.Nbytes, packet.IPAddr, pinger.Addr(), packet.Seq, packet.Rtt)
	}
	pinger.OnFinish = func(statistics *Statistics) {
		fmt.Printf("--- %s ping statistics ---\n", statistics.Addr)
		fmt.Printf("%d packets transmitted, %d received, %v%% packet loss, time %v\n",
			statistics.PacketsSent, statistics.PacketsRecv, statistics.PacketsLoss, statistics.TotalRtt)
		fmt.Printf("rtt min/avg/max/mdev = %v/%v/%v/%v \n",
			statistics.MinRtt, statistics.AvgRtt, statistics.MaxRtt, statistics.StdDevRtt)
	}
	fmt.Printf("PING %s (%s) 56(84) bytes of data.\n",
		pinger.Addr(), pinger.IPAddr())
	pinger.Run()
}


