package dns

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net"

	"golang.org/x/net/dns/dnsmessage"
)

const ROOT_SERVERS = "198.41.0.4,199.9.14.201,192.33.4.12,199.7.91.13,192.203.230.10,192.5.5.241,192.112.36.4,198.97.190.53"

func handlePacket(pc net.PacketConn, addr net.Addr, buf []byte) error {
	return fmt.Errorf("not implemented yet")
}

func outgoingDnsQuery(servers []net.IP, question dnsmessage.Question) (*dnsmessage.Parser, *dnsmessage.Header, error) {
	max := uint16(^uint16(0))

	// generate a random number max to unit16
	randomNumber, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return nil, nil, err
	}

	// convert random number type from bigint to uint16
	id := uint16(randomNumber.Int64())

	// define a UDP message (question)
	msg := dnsmessage.Message{
		Header: dnsmessage.Header{
			ID: id,
			Response: false,
			OpCode: dnsmessage.OpCode(0),
		},
		Questions: []dnsmessage.Question{question},
	}
	
	return nil, nil, nil
}