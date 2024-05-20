package packet

import (
	"fmt"
	"net"
	"time"

	"github.com/rs/zerolog/log"
)

type PacketV1 struct {
	header     []byte
	payload    []byte
	retryCount int
	retryDelay time.Duration
}

func NewV1(header []byte, payload []byte) *PacketV1 {
	return &PacketV1{
		header:  header,
		payload: payload,
	}
}

func (p *PacketV1) SetRetryCount(count int) {
	p.retryCount = count
}

func (p *PacketV1) SetRetryDelay(delay time.Duration) {
	p.retryDelay = delay
}

// Write writes packet into connection. Serves both options with delay and without
func (p *PacketV1) Write(conn net.Conn) error {
	if p.retryCount == 0 {
		return p.write(conn)
	}

	for retry := range p.retryCount {
		err := p.write(conn)
		if err == nil {
			break
		}

		if retry == p.retryCount {
			return fmt.Errorf("failed to write packet after %n retries: %w", retry, err)
		}

		log.Debug().Msg("faulty packet, sleep and retry")

		time.Sleep(p.retryDelay)
	}

	return nil
}

func (p *PacketV1) write(conn net.Conn) error {
	pkt := append(p.header, p.payload...)
	_, err := conn.Write(pkt)
	return err
}
