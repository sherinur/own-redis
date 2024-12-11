package main

import (
	"flag"
	"os"
	"strconv"

	. "own-redis/internal"
)

var portInt int

func init() {
	flag.IntVar(&portInt, "port", 8080, "Port number.")
	flag.Usage = CustomUsage
}

func main() {
	flag.Parse()
	port := "0.0.0.0:" + strconv.Itoa(portInt)

	cfg := &Config{Port: port}

	server := NewServer(cfg)

	err := server.Start()
	if err != nil {
		os.Exit(1)
	}
	// conn, err := net.ListenPacket("udp", port)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// defer conn.Close()

	// for {
	// 	buffer := make([]byte, 1024)

	// 	n, addr, err := conn.ReadFrom(buffer)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		continue
	// 	}

	// 	fmt.Println(strings.ToLower(string(buffer[:n])))
	// 	if strings.ToLower(string(buffer[:n])) == "ping\n" {
	// 		_, err = conn.WriteTo([]byte("PONG\n"), addr)
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	// 	}
	// }
}
