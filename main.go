package main

import (
	"net"

	"github.com/kehiy/berjis/core"
	"github.com/kehiy/berjis/logger"
)

func main() {
	pc, err := net.ListenPacket("udp", ":53")
	if err != nil {
		logger.Panic("can't run the DNS server:", "error", err)
	}
	defer pc.Close()

	for {
		buf := make([]byte, 512)
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			logger.Info("read error", "from", addr.String(), "network", addr.Network(), "error", err)
			continue
		}
		go core.HandlePacketOut(pc, addr, buf[:n])
	}
}
