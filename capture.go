package main

import (
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func capture(config CaptureConfig) {
	// Open device
	handle, err := pcap.OpenLive(config.Interface, 100, true, pcap.BlockForever)
	if err != nil {
		panic(err)
	}
	defer handle.Close()

	handle.SetBPFFilter(config.Filter)

	// Use the handle as a packet source to process all packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		handlePacket(packet)
	}
}

func handlePacket(packet gopacket.Packet) {
	var key string

	if layer := packet.Layer(layers.LayerTypeARP); layer != nil {
		key = "ARP"
	} else if layer := packet.Layer(layers.LayerTypeICMPv6); layer != nil {
		icmpLayer := layer.(*layers.ICMPv6)
		key = fmt.Sprintf("ICMP6-%d", icmpLayer.TypeCode.Type())
	} else if layer := packet.Layer(layers.LayerTypeIGMP); layer != nil {
		key = "IGMP"
	} else if layer := packet.Layer(layers.LayerTypeUDP); layer != nil {
		udpLayer := layer.(*layers.UDP)

		if ip := packet.Layer(layers.LayerTypeIPv6); ip != nil {
			key = fmt.Sprintf("UDP6-%d", udpLayer.DstPort)
			//fmt.Printf("UDPv6 Dst=%s Port=%d\n", ip.(*layers.IPv6).DstIP.String(), udpLayer.DstPort)
		} else if ip := packet.Layer(layers.LayerTypeIPv4); ip != nil {
			key = fmt.Sprintf("UDP4-%d", udpLayer.DstPort)
			//fmt.Printf("UDPv4 Dst=%s Port=%d\n", ip.(*layers.IPv4).DstIP.String(), udpLayer.DstPort)
		}
	} else if layer := packet.Layer(layers.LayerTypeLLC); layer != nil {
		key = "LLC"
	} else {
		key = "other"
		if layer := packet.Layer(layers.LayerTypeEthernet); layer != nil {
			ethLayer := layer.(*layers.Ethernet)
			if ethLayer.EthernetType == 0x88bf {
				key = "STP"
			} else {
				fmt.Printf("Ethernet type=%d\n", ethLayer.EthernetType)
			}
		} else {
			fmt.Printf("%s\n", packet)
		}
	}

	stats.add(key, uint16(packet.Metadata().Length))
}
