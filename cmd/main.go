package main

import (
	"flag"
	"fmt"
	"net"
	"time"

	"toydns"
)

func main() {
	fmt.Println("Start ToyDNS Server")

	addr := flag.String("resolver", "", "resolver server")
	flag.Parse()
	resolverAddr, err := net.ResolveUDPAddr("udp", *addr)
	if err != nil {
		fmt.Println("Failed to resolve resolver UDP address:", err)
		return
	}

	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2025")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	// Buffer
	udpConn.SetReadBuffer(65535)
	udpConn.SetWriteBuffer(65535)

	// Deadline
	udpConn.SetReadDeadline(time.Now().Add(5 * time.Second))

	for {

		buf := make([]byte, 512)
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				udpConn.SetReadDeadline(time.Now().Add(5 * time.Second))
				continue
			}
			fmt.Println("ReadFromUDP error:", err)
			break
		}

		fmt.Printf("received %d byte data from %s\n", size, source)
		go toydns.HandleDNS(udpConn, buf[:size], source, resolverAddr)
	}
}
