package toydns

import (
	"encoding/binary"
	"fmt"
)

// 12 bytes 0-95 bit
type Header struct {
	//0-15 bit
	Id uint16

	// 16-23 bit
	QueryResponse       bool
	Opcode              uint8
	AuthoritativeAnswer bool
	Truncation          bool
	RecursionDesired    bool

	// 24-31 bit
	RecursionAvailable bool
	Reserved           uint8
	ResponseCode       uint8

	// 32-47 bit
	QuestionCount uint16

	// 48-63 bit
	AnswerRecordCount uint16

	// 64-79 bit
	AuthorityRecordCount uint16

	// 80-95 bit
	AdditionalRecordCount uint16
}

func (h *Header) ToBytes() []byte {

	bytes := make([]byte, 12)

	// 0-15
	bytes[0] = byte(h.Id >> 8)
	bytes[1] = byte((h.Id << 8) >> 8)

	//16-23 bit
	bytes[2] = (map[bool]byte{false: 0, true: 1}[h.RecursionDesired])
	bytes[2] = bytes[2] | (map[bool]byte{false: 0, true: 1}[h.Truncation] << 1)
	bytes[2] = bytes[2] | (map[bool]byte{false: 0, true: 1}[h.AuthoritativeAnswer] << 2)
	bytes[2] = bytes[2] | (h.Opcode << 3)
	bytes[2] = bytes[2] | (map[bool]byte{false: 0, true: 1}[h.QueryResponse] << 7)

	// 24-31 bit
	bytes[3] = (h.ResponseCode) | ((h.Reserved) << 4)
	bytes[3] = bytes[3] | (map[bool]byte{false: 0, true: 1}[h.RecursionAvailable] << 7)

	// 32-47 bits (QuestionCount)
	binary.BigEndian.PutUint16(bytes[4:], h.QuestionCount)

	// 48-63 bits (AnswerRecordCount)
	binary.BigEndian.PutUint16(bytes[6:], h.AnswerRecordCount)

	// 64-79 bits (AuthorityRecordCount)
	binary.BigEndian.PutUint16(bytes[8:], h.AuthorityRecordCount)

	// 80-95 bits (AdditionalRecordCount)
	binary.BigEndian.PutUint16(bytes[10:], h.AdditionalRecordCount)

	return bytes
}

func ParseHeader(buf []byte) (*Header, error) {
	buffer := make([]byte, len(buf))
	copy(buffer, buf)

	if len(buffer) < 12 {
		return nil, fmt.Errorf("buffer is small")
	}
	id := uint16(buffer[0])<<8 | uint16(buffer[1])

	qr := (buffer[2] >> 7) == 1
	opcode := uint8((buffer[2] << 1) >> 4)
	aa := ((buffer[2] << 5) >> 7) == 1
	tc := ((buffer[2] << 6) >> 7) == 1
	rd := ((buffer[2] << 7) >> 7) == 1
	ra := (buffer[3] >> 7) == 1
	reserved := uint8((buffer[3] << 1) >> 5)
	rcode := uint8((buffer[3] << 4) >> 4)
	qcount := uint16(buffer[4])<<8 | uint16(buffer[5])
	arcount := uint16(buffer[6])<<8 | uint16(buffer[7])
	authoritycount := uint16(buffer[8])<<8 | uint16(buffer[9])
	additionalcount := uint16(buffer[10])<<8 | uint16(buffer[11])
	return &Header{
		Id:                    id,
		QueryResponse:         qr,
		Opcode:                opcode,
		AuthoritativeAnswer:   aa,
		Truncation:            tc,
		RecursionDesired:      rd,
		RecursionAvailable:    ra,
		Reserved:              reserved,
		ResponseCode:          rcode,
		QuestionCount:         qcount,
		AnswerRecordCount:     arcount,
		AuthorityRecordCount:  authoritycount,
		AdditionalRecordCount: additionalcount,
	}, nil
}
