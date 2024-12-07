package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	. "own-redis/internal"
)

var portInt int

func init() {
	flag.IntVar(&portInt, "port", 8080, "Port number.")
	flag.Usage = CustomUsage
}

func main() {
	flag.Parse()

	port := ":" + strconv.Itoa(portInt)

	conn, err := net.ListenPacket("udp", port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("Own Redis is started and listening on %s\n", port[1:])

	for {
		buffer := make([]byte, 1024)

		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(strings.ToLower(string(buffer[:n])))
		if strings.ToLower(string(buffer[:n])) == "ping\n" {
			_, err = conn.WriteTo([]byte("PONG\n"), addr)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
