package toydns

import (
	"fmt"
	"log"
	"net"
	"time"
)

func HandleDNS(udpConn *net.UDPConn, req []byte, source *net.UDPAddr, resolverAddr *net.UDPAddr) {
	response := []byte{}
	header, err := ParseHeader(req)
	if err != nil {
		log.Fatalf("parse header error %v", err)
	}
	Id := header.Id

	if header.Opcode != 0 {
		header.Id = Id
		header.ResponseCode = 4
		header.QueryResponse = true
		header.AuthoritativeAnswer = false
		header.Truncation = false
		header.RecursionAvailable = false
		header.Reserved = 0
		header.AuthorityRecordCount = 0
		header.AdditionalRecordCount = 0

		header_bytes := header.ToBytes()
		response = append(response, header_bytes...)
		response = append(response, req[12:]...)
		_, err = udpConn.WriteToUDP(response, source)

		if err != nil {
			fmt.Println("failed to send response:", err)
		}
		return
	}

	forward_response, err := Forward(req, resolverAddr)
	if err != nil {
		log.Fatalf("forward error %v", err)
	}

	resp_header, err := ParseHeader(forward_response)
	if err != nil {
		log.Fatalf("parse header error %v", err)
	}
	resp_header.Id = Id
	resp_header.AuthoritativeAnswer = false
	resp_header.Reserved = 0
	resp_header.Truncation = false
	resp_header.RecursionAvailable = false
	resp_header.Reserved = 0
	resp_header.AuthorityRecordCount = 0
	resp_header.AdditionalRecordCount = 0

	response = append(resp_header.ToBytes(), forward_response[12:]...)

	_, err = udpConn.WriteToUDP(response, source)
	if err != nil {
		fmt.Println("failed to send response:", err)
	}
}

func Forward(req []byte, resolverAddr *net.UDPAddr) ([]byte, error) {
	conn, err := net.Dial("udp", resolverAddr.String())
	if err != nil {
		return nil, fmt.Errorf("dial fail error: %v", err)
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(5 * time.Second))

	_, err = conn.Write(req)
	if err != nil {
		return nil, fmt.Errorf("send request error: %v", err)
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("receive fail!", err)
		return nil, fmt.Errorf("receive fail for request %v", err)
	}
	return buf[:n], nil
}
