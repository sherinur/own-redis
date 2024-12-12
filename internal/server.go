package internal

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	config  *Config
	storage *storage

	conn     *net.UDPConn
	stopChan chan os.Signal // chan to handle signals of exit
	doneChan chan struct{}  // chan to sync goroutines
}

func NewServer(cfg *Config) *Server {
	return &Server{
		config:   cfg,
		storage:  NewStorage(),
		stopChan: make(chan os.Signal, 1),
		doneChan: make(chan struct{}),
	}
}

type Config struct {
	Port string
}

func (s *Server) Start() error {
	signal.Notify(s.stopChan, syscall.SIGINT, syscall.SIGTERM)

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
	s.conn = conn
	defer conn.Close()

	slog.Info(fmt.Sprintf("Own Redis is started and listening on address %s", s.config.Port))

	go func() {
		<-s.stopChan
		slog.Info("Shutdown signal received. Stopping server...")

		s.Shutdown()
		close(s.doneChan)
	}()

	for {
		buff := make([]byte, 1024)

		n, addr, err := conn.ReadFromUDP(buff)
		if err != nil {
			slog.Error(fmt.Sprintf("Error of reading from UDP connection: %s", err.Error()))
		}

		slog.Info(fmt.Sprintf("Received %d bytes from %s: %s", n, addr.IP, string(buff[:n])))

		if n == 0 || string(buff) == "\n" {
			continue
		}

		// handling request
		go func() {
			if err := handleUDPRequest(s, conn, addr, buff); err != nil {
				fmt.Println(err.Error())
			}
		}()
	}

	<-s.doneChan
	return nil
}

func (s *Server) Shutdown() {
	slog.Info("Starting gracefull shutdown...")

	// closing udp conn
	if err := s.conn.Close(); err != nil {
		slog.Error(fmt.Sprintf("Error of closing UDP connection: %s", err.Error()))
	} else {
		slog.Info("UDP connection closed.")
	}

	// wait for other goroutines
	<-s.doneChan

	slog.Info("Server gracefully shut down.")
}
