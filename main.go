package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// header:
// - protocol version
// - packet type
// body

const (
	protocolVersion1  uint8 = 1
	packetTypePing    uint8 = 10
	packetTypePong    uint8 = 11
	packetTypePayload uint8 = 20
)

func main() {
	setLogger()

	listerAddress := os.Getenv("LISTEN_ADDRESS")

	l, err := net.Listen("tcp", listerAddress)
	if err != nil {
		log.Fatal().Err(err).Str("address", listerAddress).Msg("failed to start listener")
	}

	defer func() {
		if err := l.Close(); err != nil {
			log.Error().Err(err).Str("address", listerAddress).Msg("failed to close listener")
		}
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Error().Err(err).Msg("failed to accept connection")
			continue
		}

		go handleConnection(conn)
	}
}

func writePong(conn net.Conn) error {
	header := make([]byte, 4)
	header[0] = protocolVersion1
	header[1] = packetTypePong

	body := make([]byte, 10)

	packet := append(header, body...)

	_, err := conn.Write(packet)
	if err != nil {
		return fmt.Errorf("failed to write packet: %w", err)
	}

	return nil
}

func setLogger() {
	level, exists := os.LookupEnv("LOG_LEVEL")
	if !exists {
		level = "info"
	}

	zeroLvl, err := zerolog.ParseLevel(level)
	if err != nil {
		zeroLvl = zerolog.InfoLevel
		log.Info().Err(err).Msg("failed to parse log level, using info")
	}

	zerolog.SetGlobalLevel(zeroLvl)
	log.Info().Str("set_level", zeroLvl.String()).Msg("logger is set")
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		buf := make([]byte, 14)
		n, err := io.ReadFull(conn, buf)
		if err != nil {
			if n == 0 && errors.Is(err, io.EOF) {
				continue
			}

			log.Error().Err(err).Msg("failed to read")
			continue
		}

		if len(buf) < 14 {
			log.Error().Msg("packet len is below required")
			continue
		}

		switch buf[1] {
		case packetTypePing:
			log.Debug().Msg("ping")
			if err = writePong(conn); err != nil {
				log.Error().Err(err).Msg("failed to write pong packet")
			}
			log.Debug().Msg("pong")
		case 20:
		default:
			log.Error().Msg("unknown packet type")
		}
	}
}
