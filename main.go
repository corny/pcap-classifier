package main

import (
	"log"

	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	device = "wlp4s0"
	stats  = make(Stats)
)

// Counter counts bytes and packets per packet type
type Counter struct {
	Bytes   uint64
	Packets uint
}

// Stats is a map from packet type to counter
type Stats map[string]*Counter

/*
[proto: "UDPv4", type: "5353"]
[proto: "UDPv4", type: "1234"]
[proto: "IGMP"]
*/

func main() {
	// Open device
	handle, err := pcap.OpenLive(device, 1500, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	handle.SetBPFFilter("ether multicast")

	// Use the handle as a packet source to process all packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		// Process packet here
		if layer := packet.Layer(layers.LayerTypeARP); layer != nil {
			stats.add("ARP", packet)
		} else if layer := packet.Layer(layers.LayerTypeICMPv6); layer != nil {
			icmpLayer := layer.(*layers.ICMPv6)
			stats.add(fmt.Sprintf("ICMPv6-%d", icmpLayer.TypeCode), packet)
		} else if layer := packet.Layer(layers.LayerTypeIGMP); layer != nil {
			stats.add("IGMP", packet)
		} else if layer := packet.Layer(layers.LayerTypeUDP); layer != nil {
			udpLayer := layer.(*layers.UDP)

			if ip := packet.Layer(layers.LayerTypeIPv6); ip != nil {
				stats.add(fmt.Sprintf("UDPv6-%d", udpLayer.DstPort), packet)
				//fmt.Printf("UDPv6 Dst=%s Port=%d\n", ip.(*layers.IPv6).DstIP.String(), udpLayer.DstPort)
			} else if ip := packet.Layer(layers.LayerTypeIPv4); ip != nil {
				stats.add(fmt.Sprintf("UDPv4-%d", udpLayer.DstPort), packet)
				//fmt.Printf("UDPv4 Dst=%s Port=%d\n", ip.(*layers.IPv4).DstIP.String(), udpLayer.DstPort)
			}
		} else {
			fmt.Printf("%s\n", packet)
		}
	}
}

func (stats Stats) add(key string, packet gopacket.Packet) {
	counter, _ := stats[key]

	if counter == nil {
		counter = &Counter{}
	}

	log.Println(key, len(packet.Data()))

	counter.Bytes += uint64(len(packet.Data()))
	counter.Packets++
}
