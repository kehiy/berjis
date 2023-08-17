package main

import (
	"fmt"
	"log"
	"net"

	"github.com/kehiy/dns-server/core"
)

func main() {
	pc, err := net.ListenPacket("udp", ":53")
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	for {
		buf := make([]byte, 512)
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			fmt.Printf("Connection error [%s]: %s\n", addr.String(), err)
			continue
		}
		go core.HandlePacket(pc, addr, buf[:n])
	}
}
