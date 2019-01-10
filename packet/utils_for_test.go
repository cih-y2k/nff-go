// +build !test

package packet

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/intel-go/nff-go/common"
	"github.com/intel-go/nff-go/low"
)

// isInit is common for all tests
var isInit bool

const testDevelopmentMode = false
const payloadSize = 100

func tInitDPDK() {
	if isInit != true {
		argc, argv := low.InitDPDKArguments([]string{})
		// burstSize=32, mbufNumber=8191, mbufCacheSize=250
		if err := low.InitDPDK(argc, argv, 32, 8191, 250, 0); err != nil {
			log.Fatal(err)
		}
		nonPerfMempool = low.CreateMempool("Test")
		isInit = true
	}
}

func getIPv4TCPTestPacket() *Packet {
	pkt := getPacket()
	InitEmptyIPv4TCPPacket(pkt, payloadSize)

	initEtherAddrs(pkt)
	initIPv4Addrs(pkt)
	initPorts(pkt)

	return pkt
}

func getIPv4UDPTestPacket() *Packet {
	pkt := getPacket()
	InitEmptyIPv4UDPPacket(pkt, payloadSize)

	initEtherAddrs(pkt)
	initIPv4Addrs(pkt)
	initPorts(pkt)

	return pkt
}

func getIPv4ICMPTestPacket() *Packet {
	pkt := getPacket()
	InitEmptyIPv4ICMPPacket(pkt, payloadSize)

	initEtherAddrs(pkt)
	initIPv4Addrs(pkt)

	return pkt
}

func getIPv6TCPTestPacket() *Packet {
	pkt := getPacket()
	InitEmptyIPv6TCPPacket(pkt, payloadSize)

	initEtherAddrs(pkt)
	initIPv6Addrs(pkt)
	initPorts(pkt)

	return pkt
}

func getIPv6UDPTestPacket() *Packet {
	pkt := getPacket()
	InitEmptyIPv6UDPPacket(pkt, payloadSize)

	initEtherAddrs(pkt)
	initIPv6Addrs(pkt)
	initPorts(pkt)
	return pkt
}

func getIPv6ICMPTestPacket() *Packet {
	pkt := getPacket()
	InitEmptyIPv6ICMPPacket(pkt, payloadSize)
	initEtherAddrs(pkt)
	initIPv6Addrs(pkt)

	return pkt
}

func getARPRequestTestPacket() *Packet {
	pkt := getPacket()

	sha := common.MACAddress{0x01, 0x11, 0x21, 0x31, 0x41, 0x51}
	spa := common.SliceToIPv4(net.ParseIP("127.0.0.1").To4())
	tpa := common.SliceToIPv4(net.ParseIP("128.9.9.5").To4())
	InitARPRequestPacket(pkt, sha, spa, tpa)

	return pkt
}

func initEtherAddrs(pkt *Packet) {
	pkt.Ether.SAddr = common.MACAddress{0x01, 0x11, 0x21, 0x31, 0x41, 0x51}
	pkt.Ether.DAddr = common.MACAddress{0x0, 0x11, 0x22, 0x33, 0x44, 0x55}
}

func initIPv4Addrs(pkt *Packet) {
	pkt.GetIPv4().SrcAddr = common.SliceToIPv4(net.ParseIP("127.0.0.1").To4())
	pkt.GetIPv4().DstAddr = common.SliceToIPv4(net.ParseIP("128.9.9.5").To4())
}

func initIPv6Addrs(pkt *Packet) {
	copy(pkt.GetIPv6().SrcAddr[:], net.ParseIP("dead::beaf")[:common.IPv6AddrLen])
	copy(pkt.GetIPv6().DstAddr[:], net.ParseIP("dead::beaf")[:common.IPv6AddrLen])
}

func initPorts(pkt *Packet) {
	// Src and Dst port numbers placed at the same offset from L4 start in both tcp and udp
	l4 := (*UDPHdr)(pkt.L4)
	l4.SrcPort = SwapBytesUint16(1234)
	l4.DstPort = SwapBytesUint16(5678)
}

func getPacket() *Packet {
	pkt, err := NewPacket()
	if err != nil {
		log.Fatal(err)
	}
	return pkt
}

func dumpPacketToPcap(fileName string, pkt *Packet) {
	if !testDevelopmentMode {
		return
	}

	file, err := os.Create(fileName + ".pcap")
	if err != nil {
		fmt.Println(err)
	}
	err = WritePcapGlobalHdr(file)
	if err != nil {
		fmt.Println(err)
	}
	err = pkt.WritePcapOnePacket(file)
	if err != nil {
		fmt.Println(err)
	}
	err = file.Close()
	if err != nil {
		fmt.Println(err)
	}
}
