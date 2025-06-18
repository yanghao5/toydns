package toydns

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

type Question struct {
	QName  string
	QType  uint16
	QClass uint16
}

func (q Question) ToBytes() []byte {
	var result []byte
	result = append(result, EncodeDomain(q.QName)...)

	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, q.QType)
	binary.Write(&buf, binary.BigEndian, q.QClass)

	result = append(result, buf.Bytes()...)
	return result
}

func ParseQuestion(buffer []byte) (qusetion *Question, offset int, err error) {
	if len(buffer) < 13 {
		return nil, 12, fmt.Errorf("data too small")
	}
	offset = 12

	qname, move := DecodeDomain(buffer[offset:])
	offset = offset + move
	qtype := binary.BigEndian.Uint16(buffer[offset : offset+2])
	qclass := binary.BigEndian.Uint16(buffer[offset+2 : offset+4])

	offset += 4
	return &Question{
		QName:  qname,
		QType:  qtype,
		QClass: qclass,
	}, offset, nil
}

func EncodeDomain(domain string) []byte {
	var result []byte
	labels := strings.Split(domain, ".")
	for _, label := range labels {
		result = append(result, byte(len(label)))
		result = append(result, []byte(label)...)
	}
	result = append(result, 0x00)
	return result
}

// DecodeDomain parses a domain name from the DNS message buffer starting at domain encoding position.
// It returns the parsed domain name and the number of bytes that were read.
//
// Example:
//
//	buf := []byte{0,0,0,0,0,3, 'w', 'w', 'w', 6, 'g', 'o', 'o', 'g', 'l', 'e', 3, 'c', 'o', 'm', 0}
//	domain, move := DecodeDomain(buf[6:])
//	// domain == "www.google.com"
//	// move == 17
func DecodeDomain(buf []byte) (domain string, move int) {
	var result []string
	i := 0

	for {
		length := buf[i]
		// label index
		i++
		if length == 0 {
			break
		}
		label := string(buf[i : i+int(length)])
		result = append(result, label)

		// next label
		i += int(length)
	}

	domain = strings.Join(result, ".")
	return domain, i
}
