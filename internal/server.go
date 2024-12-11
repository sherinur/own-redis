package internal

import (
	"fmt"
	"log/slog"
	"net"
)

type Server struct {
	config  *Config
	storage *storage
}

func NewServer(cfg *Config) *Server {
	return &Server{
		config:  cfg,
		storage: NewStorage(),
	}
}

type Config struct {
	Port string
}

func (s *Server) Start() error {
	udpAddr, err := net.ResolveUDPAddr("udp", s.config.Port)
	if err != nil {
		slog.Error(fmt.Sprintf("Error of resolving UDPAddr: %s", err.Error()))
		return err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		slog.Error(fmt.Sprintf("Error of starting UDP server: %s", err.Error()))
		return err
	}
	defer conn.Close()

	slog.Info(fmt.Sprintf("Own Redis is started and listening on port %s", s.config.Port))

	for {
		buff := make([]byte, 1024)

		n, addr, err := conn.ReadFromUDP(buff)
		if err != nil {
			slog.Error(fmt.Sprintf("Error of reading from UDP connection: %s", err.Error()))
		}

		slog.Info(fmt.Sprintf("Received %d bytes from %s: %s", n, addr.IP, string(buff[:n-1])))

		if n == 0 || string(buff) == "\n" {
			continue
		}

		err = handleUDPRequest(s, conn, addr, buff)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}
}
