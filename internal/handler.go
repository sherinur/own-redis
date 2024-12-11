package internal

import (
	"fmt"
	"net"
	"strings"
)

func handleUDPRequest(s *Server, conn *net.UDPConn, addr *net.UDPAddr, buff []byte) error {
	input := strings.TrimSpace(string(buff))

	fields := strings.Fields(input)
	if len(fields) == 0 {
		writeErrorResponse(conn, addr, "empty command")
		return nil
	}
	fields = fields[:len(fields)-1]

	cmd := fields[0]

	switch strings.ToUpper(cmd) {
	case "PING":
		conn.WriteToUDP([]byte("PONG\n"), addr)
	case "SET":
		if len(fields) < 3 {
			writeErrorResponse(conn, addr, "wrong number of arguments for 'set' command")
			break
		}

		key, value := fields[1], fields[2]
		err := s.storage.Set(key, value)
		if err != nil {
			writeErrorResponse(conn, addr, "OOM command not allowed when used memory > 'maxmemory'")
		}

		writeUDPResponse(conn, addr, "OK")
	case "GET":
		if len(fields) < 2 {
			writeErrorResponse(conn, addr, "wrong number of arguments for 'get' command")
			break
		}

		key := fields[1]
		value := s.storage.Get(key)

		if value == nil {
			writeUDPResponse(conn, addr, "(nil)")
		} else {
			writeUDPResponse(conn, addr, *value)
		}
	default:
		response := fmt.Sprintf("unknown command '%s'", cmd)
		writeErrorResponse(conn, addr, response)
	}

	return nil
}

func writeUDPResponse(conn *net.UDPConn, addr *net.UDPAddr, msg string) {
	conn.WriteToUDP([]byte(msg+"\n"), addr)
}

func writeErrorResponse(conn *net.UDPConn, addr *net.UDPAddr, msg string) {
	response := fmt.Sprintf("(error) ERR %s", msg)
	writeUDPResponse(conn, addr, response)
}
