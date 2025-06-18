package toydns

import (
	"fmt"
	"net"
)

func HandleDNS(udpConn *net.UDPConn, req []byte, source *net.UDPAddr, resolverAddr *net.UDPAddr) {
	response := req
	_ = resolverAddr

	_, err := udpConn.WriteToUDP(response, source)
	if err != nil {
		fmt.Println("failed to send response:", err)
	}
}
