package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/sergkondr/dinosaur/internal/dns"
)

func main() {
	dnsDatabase := map[string][][]int{
		"test.com":         {{127, 0, 0, 1}},
		"test.example.com": {{127, 0, 0, 2}, {127, 0, 0, 3}},
	}

	listenAddr := net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 9053,
	}

	connection, err := net.ListenUDP("udp", &listenAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Dinosaur is listening on", listenAddr.String())
	for {
		var packet [512]byte

		n, addr, err := connection.ReadFromUDP(packet[0:])
		if err != nil {
			fmt.Println("Error reading from UDP socket:", err)
		}

		if n != 0 {
			msg, err := dns.Deserialize(packet[0:])
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(msg.String())

			dnsName := strings.Join(msg.Question.Qname, ".")
			addresses, ok := dnsDatabase[dnsName]
			if !ok {
				fmt.Println("not found any addresses for ", dnsName)
			} else {
				fmt.Println("found following addresses for ", msg.Question.Qname, ":", addresses)
				resp := createResponseDNSPacket(msg.Question.Qname, addresses, msg.Header.ID)
				respData, err := resp.Serialize()
				if err != nil {
					fmt.Println(err)
				}

				_, err = connection.WriteToUDP(respData, addr)
				if err != nil {
					fmt.Println("Error writing to UDP socket:", err)
				}
			}
		}
	}
}

func createResponseDNSPacket(dnsName []string, addresses [][]int, id uint16) dns.Message {
	fmt.Println("creating response for ", dnsName, " with ID", id, " and addresses", addresses)

	msg := dns.Message{}
	msg.Header.ID = id
	msg.Header.Flags = dns.Flags{
		Query:               1,
		Opcode:              0,
		AuthoritativeAnswer: 0,
		Truncation:          0,
		RecursionDesired:    1,
		RecursionAvailable:  1,
		Z:                   0,
		ResponseCode:        0,
	}

	msg.Question = dns.MessageQuestion{
		Qname:  dnsName,
		Qclass: dns.RECORD_CLASS_IN,
		Qtype:  dns.RECORD_TYPE_A,
	}

	for i := 0; i < len(addresses); i++ {
		msg.Answer = append(msg.Answer, dns.MessageResourceRecord{
			Name:  dnsName,
			Type:  dns.RECORD_TYPE_A,   // https://datatracker.ietf.org/doc/html/rfc1035#section-3.2.2
			Class: dns.RECORD_CLASS_IN, // 1 = Internet
			TTL:   300,
			Data:  addresses[i],
		})
	}

	return msg
}
