package main

import (
	"net"

	"github.com/kehiy/berjis/core"
	"github.com/kehiy/berjis/log"
	"github.com/spf13/cobra"
)

func buildRunCommand(parentCmd *cobra.Command) {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "run a berjis instance",
	}
	parentCmd.AddCommand(runCmd)

	runCmd.Run = func(_ *cobra.Command, _ []string) {
		pc, err := net.ListenPacket("udp", ":53")
		if err != nil {
			log.Panic("can't run the DNS server:", "error", err)
		}
		defer pc.Close()

		for {
			buf := make([]byte, core.MAXPACKETSIZE)
			n, addr, err := pc.ReadFrom(buf)
			if err != nil {
				log.Info("read error", "from", addr.String(), "network", addr.Network(), "error", err)
				continue
			}
			go core.HandlePacketOut(pc, addr, buf[:n])
		}
	}
}
