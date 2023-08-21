package core

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"strings"

	"github.com/kehiy/berjis/logger"
	"golang.org/x/net/dns/dnsmessage"
)

const ROOTSERVERS = `198.41.0.4,199.9.14.201,192.33.4.12,199.7.91.13,
					 192.203.230.10,192.5.5.241,192.112.36.4,198.97.190.53`

func HandlePacketOut(pc net.PacketConn, addr net.Addr, buf []byte) {
	//* send incoming packets to handlePacket function
	if err := HandlePacketIn(pc, addr, buf); err != nil {
		fmt.Printf("handlePacket error [%s]: %s\n", addr.String(), err)
	}
}

func HandlePacketIn(pc net.PacketConn, addr net.Addr, buf []byte) error {
	p := dnsmessage.Parser{}
	header, err := p.Start(buf)
	if err != nil {
		return err
	}

	question, err := p.Question()
	if err != nil {
		return err
	}

	response, err := DNSQuery(getRootServers(), question)
	if err != nil {
		return err
	}
	response.Header.ID = header.ID

	responseBuffer, err := response.Pack()
	if err != nil {
		return err
	}

	//* send response
	_, err = pc.WriteTo(responseBuffer, addr)
	if err != nil {
		return err
	}

	return nil
}

func DNSQuery(servers []net.IP, question dnsmessage.Question) (*dnsmessage.Message, error) {
	for i := 0; i < 3; i++ {
		dnsAnswer, header, err := OutgoingDNSQuery(servers, question)
		if err != nil {
			return nil, err
		}

		parsedAnswers, err := dnsAnswer.AllAnswers()
		if err != nil {
			return nil, err
		}

		//! if the header is Authoritative we send the ip's back here.
		if header.Authoritative {
			return &dnsmessage.Message{
				Header:  dnsmessage.Header{Response: true},
				Answers: parsedAnswers,
			}, nil
		}

		authorities, err := dnsAnswer.AllAuthorities()
		if err != nil {
			return nil, err
		}

		//* if there is no authority we send a server failure code
		if len(authorities) == 0 {
			return &dnsmessage.Message{
				Header: dnsmessage.Header{RCode: dnsmessage.RCodeNameError},
			}, nil
		}

		//* get name servers in the response
		nameservers := make([]string, len(authorities))
		for k, authority := range authorities {
			if authority.Header.Type == dnsmessage.TypeNS {
				nameservers[k] = authority.Body.(*dnsmessage.NSResource).NS.String()
			}
		}

		additionals, err := dnsAnswer.AllAdditionals()
		if err != nil {
			return nil, err
		}

		newResolverServersFound := false

		//* change servers list to new name servers get from previous server...
		servers = []net.IP{}
		for _, additional := range additionals {
			if additional.Header.Type == dnsmessage.TypeA {
				for _, nameserver := range nameservers {
					if additional.Header.Name.String() == nameserver {
						newResolverServersFound = true
						servers = append(servers, additional.Body.(*dnsmessage.AResource).A[:])
					}
				}
			}
		}

		if !newResolverServersFound {
			for _, nameserver := range nameservers {
				if !newResolverServersFound {
					response, err := DNSQuery(getRootServers(),
						dnsmessage.Question{
							Name: dnsmessage.MustNewName(nameserver),
							Type: dnsmessage.TypeA, Class: dnsmessage.ClassINET,
						})
					if err != nil {
						logger.Warn("lookup of nameserver failed", "nameserver", nameserver, "error", err)
					} else {
						newResolverServersFound = true
						for _, answer := range response.Answers {
							if answer.Header.Type == dnsmessage.TypeA {
								servers = append(servers, answer.Body.(*dnsmessage.AResource).A[:])
							}
						}
					}
				}
			}
		}
	}
	return &dnsmessage.Message{
		Header: dnsmessage.Header{RCode: dnsmessage.RCodeServerFailure},
	}, nil
}

func OutgoingDNSQuery(servers []net.IP, question dnsmessage.Question) (*dnsmessage.Parser, *dnsmessage.Header, error) {
	max := ^uint16(0)

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
			ID:       id,
			Response: false,
			OpCode:   dnsmessage.OpCode(0),
		},
		Questions: []dnsmessage.Question{question},
	}

	// pack message to send to root server
	buf, err := msg.Pack()
	if err != nil {
		return nil, nil, err
	}

	// make connection with servers and send message (for each servers)
	//* connections is with udp protocol
	var conn net.Conn
	for _, server := range servers {
		conn, err = net.Dial("udp", server.String()+":53" /*dns port*/)
		if err == nil {
			break
		}
	}
	if conn == nil {
		return nil, nil, fmt.Errorf("faild to make connection to servers: %w", err)
	}
	// here we have the connection!
	_, err = conn.Write(buf)
	if err != nil {
		return nil, nil, err
	}

	// read answer  from connection
	answer := make([]byte, 512)
	n, err := bufio.NewReader(conn).Read(answer)
	if err != nil {
		return nil, nil, err
	}

	conn.Close()

	var p dnsmessage.Parser
	header, err := p.Start(answer[:n])
	if err != nil {
		return nil, nil, fmt.Errorf("parser start error: %w", err)
	}

	questions, err := p.AllQuestions()
	if err != nil {
		return nil, nil, err
	}
	if len(questions) != len(msg.Questions) {
		return nil, nil, fmt.Errorf("answer packet dosen't have the same amount of questions")
	}

	err = p.SkipAllQuestions()
	if err != nil {
		return nil, nil, err
	}

	return &p, &header, nil
}

// * make a loop over ROOT SERVERS list and return a slice of root servers ip.
func getRootServers() []net.IP {
	rootServers := []net.IP{}
	for _, rootServer := range strings.Split(ROOTSERVERS, ",") {
		rootServers = append(rootServers, net.ParseIP(rootServer))
	}
	return rootServers
}
