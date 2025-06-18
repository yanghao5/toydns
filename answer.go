package toydns

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

type Answer struct {
	QName      string
	QType      uint16
	QClass     uint16
	TTL        uint32
	DataLength uint16
	Data       string
}

func ipv4str_to_big_endian_bytes(ip_str string) []byte {
	ip := net.ParseIP(ip_str)
	if ip == nil {
		fmt.Println("invalid ip")
		return nil
	}
	ipv4 := ip.To4()
	if ipv4 == nil {
		fmt.Println("not IPv4")
		return nil
	}
	return ipv4
}

func bytes_to_ipv4str(ip_bytes []byte) string {
	if len(ip_bytes) != 4 {
		return ""
	}
	return fmt.Sprintf("%d.%d.%d.%d", ip_bytes[0], ip_bytes[1], ip_bytes[2], ip_bytes[3])
}

func DecodeCompressedDomain(data []byte, buf []byte) (string, int) {
	i := 0
	result := []string{}
	length := 0
	for {
		if (data[i] & 0xC0) == 0xC0 {
			pointer := int(uint16(data[i])<<8|uint16(data[i+1])) & 0x3FFF
			domain_compressed, _ := DecodeDomain(buf[pointer:])
			result = append(result, domain_compressed)
			i += 2
			break
		} else {
			length = int(data[i])
		}

		// label index
		i++
		if length == 0 {
			break
		}
		label := string(data[i : i+int(length)])
		result = append(result, label)

		// next label
		i += int(length)
	}
	domain_name := strings.Join(result, ".")
	return domain_name, i
}

func GetATypeAnswer(buffer []byte, offset int) ([]byte, int) {
	start := offset
	_, move := DecodeCompressedDomain(buffer[offset:], buffer)
	offset = offset + move

	// qtype & class & ttl
	offset += 8

	// datalength
	datalength := binary.BigEndian.Uint16(buffer[offset : offset+2])
	offset += 2

	// data
	offset += int(datalength)
	return buffer[start:offset], offset
}

//func (a Answer) ToBytes() []byte {
// 	var result []byte
// 	result = append(result, []byte{0xC0, 0x0C}...)
// 	var buf bytes.Buffer
// 	binary.Write(&buf, binary.BigEndian, a.QType)
// 	binary.Write(&buf, binary.BigEndian, a.QClass)

// 	ttl_bytes := make([]byte, 4)
// 	// binary.BigEndian.PutUint32(ttl_bytes, a.TTL)
// 	ttl_bytes[0] = byte((a.TTL >> 24) & 0xFF)
// 	ttl_bytes[1] = byte((a.TTL >> 16) & 0xFF)
// 	ttl_bytes[2] = byte((a.TTL >> 8) & 0xFF)
// 	ttl_bytes[3] = byte(a.TTL & 0xFF)
// 	buf.Write(ttl_bytes)

// 	binary.Write(&buf, binary.BigEndian, a.DataLength)

// 	buf.Write(ipv4str_to_big_endian_bytes(a.Data))
// 	result = append(result, buf.Bytes()...)
// 	return result
// }

// func ParseAnswer(buffer []byte, offset int, ancount int) ([]*Answer, int, error) {
// 	if len(buffer) < offset {
// 		return nil, offset, fmt.Errorf("data too small")
// 	}
// 	qname, move := decodeDomainName(buffer[offset:])
// 	offset = offset + move

// 	qtype := binary.BigEndian.Uint16(buffer[offset : offset+2])
// 	qclass := binary.BigEndian.Uint16(buffer[offset+2 : offset+4])
// 	offset = offset + 4

// 	ttl := binary.BigEndian.Uint32(buffer[offset : offset+4])
// 	offset += 4

// 	datalength := binary.BigEndian.Uint16(buffer[offset : offset+2])
// 	offset = offset + 2
// 	// test2
// 	data := bytes_to_ipstr(buffer[offset:datalength])
// 	offset += int(datalength)

// 	return &PlainDNSMsgAnswer{
// 		QName:      qname,
// 		QType:      qtype,
// 		QClass:     qclass,
// 		TTL:        ttl,
// 		DataLength: datalength,
// 		Data:       data,
// 	}, offset, nil
// }

// func ParseMultiPlainDNSMsgAnswer(buffer []byte, offset int) (*PlainDNSMsgAnswer, int, error) {
// 	if len(buffer) < offset {
// 		return nil, offset, fmt.Errorf("data too small")
// 	}
// 	qname, move := MultiQuestionDecodeDomainName(buffer[offset:], buffer)
// 	offset = offset + move

// 	fmt.Println(qname)
// 	fmt.Println("move ", move, "offset ", offset)

// 	qtype := binary.BigEndian.Uint16(buffer[offset : offset+2])
// 	qclass := binary.BigEndian.Uint16(buffer[offset+2 : offset+4])
// 	offset = offset + 4

// 	ttl := binary.BigEndian.Uint32(buffer[offset : offset+4])
// 	offset += 4

// 	datalength := binary.BigEndian.Uint16(buffer[offset : offset+2])
// 	offset = offset + 2
// 	// test2
// 	data := bytes_to_ipstr(buffer[offset : offset+int(datalength)])
// 	offset += int(datalength)

// 	return &PlainDNSMsgAnswer{
// 		QName:      qname,
// 		QType:      qtype,
// 		QClass:     qclass,
// 		TTL:        ttl,
// 		DataLength: datalength,
// 		Data:       data,
// 	}, offset, nil
// }
