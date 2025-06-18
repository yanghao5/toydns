package toydns_test

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"

	"toydns"
)

// go test -v -run TestHeaderToBytes
func TestHeaderToBytes(t *testing.T) {
	cases := []struct {
		name   string
		header toydns.Header
		want   []byte
	}{
		{
			name:   "All zero",
			header: toydns.Header{},
			want: []byte{
				0x00, 0x00, // Id
				0x00,       // flags byte 2
				0x00,       // flags byte 3
				0x00, 0x00, // QuestionCount
				0x00, 0x00, // AnswerRecordCount
				0x00, 0x00, // AuthorityRecordCount
				0x00, 0x00, // AdditionalRecordCount
			},
		},
		{
			name: "Full flags and counts",
			header: toydns.Header{
				Id:                    0xFFFF,
				QueryResponse:         true,
				Opcode:                0xF, // 4 bits max
				AuthoritativeAnswer:   true,
				Truncation:            true,
				RecursionDesired:      true,
				RecursionAvailable:    true,
				Reserved:              0xF, // 4 bits max
				ResponseCode:          0xF, // 4 bits max
				QuestionCount:         1,
				AnswerRecordCount:     2,
				AuthorityRecordCount:  3,
				AdditionalRecordCount: 4,
			},
			want: []byte{
				0xFF, 0xFF,
				0x01 | (1 << 1) | (1 << 2) | (0xF << 3) | (1 << 7), // flags byte 2
				0x0F | (0xF << 4) | (1 << 7),                       // flags byte 3
				0x00, 0x01,
				0x00, 0x02,
				0x00, 0x03,
				0x00, 0x04,
			},
		},
		{
			name: "case 3",
			header: toydns.Header{
				Id:                    0x1234,
				QueryResponse:         true,
				Opcode:                0x0,
				AuthoritativeAnswer:   false,
				Truncation:            false,
				RecursionDesired:      true,
				RecursionAvailable:    true,
				Reserved:              0x0,
				ResponseCode:          0x0,
				QuestionCount:         0x1,
				AnswerRecordCount:     0x6,
				AuthorityRecordCount:  0x0,
				AdditionalRecordCount: 0x0,
			},
			want: []byte{
				0x12, 0x34, // Transaction ID
				0x81, 0x80, // QR = 1 (Response), Standard Query
				0x00, 0x01, // QDCOUNT = 1
				0x00, 0x06, // ANCOUNT = 6
				0x00, 0x00, // NSCOUNT = 0
				0x00, 0x00, // ARCOUNT = 0
			},
		},
	}
	for _, c := range cases {
		got := c.header.ToBytes()
		if len(got) != len(c.want) {
			t.Errorf("%s: length mismatch got %d want %d", c.name, len(got), len(c.want))
			continue
		}
		for i := range got {
			if got[i] != c.want[i] {
				t.Errorf("%s: byte %d mismatch got %02X want %02X", c.name, i, got[i], c.want[i])
			}
		}

		t.Logf("%s: bytes literal:\n[]byte{%s}", c.name, hex.EncodeToString(got))
	}
}

// go test -v -run TestHeaderParseHeader
func TestHeaderParseHeader(t *testing.T) {
	cases := []struct {
		name   string
		header []byte
		want   toydns.Header
	}{
		{
			name: "case 1",
			header: []byte{
				0x12, 0x34, // Transaction ID
				0x81, 0x80, // Flags
				0x00, 0x01, // QDCOUNT = 1
				0x00, 0x06, // ANCOUNT = 6
				0x00, 0x00, // NSCOUNT = 0
				0x00, 0x00, // ARCOUNT = 0
			},
			want: toydns.Header{
				Id:                    0x1234,
				QueryResponse:         true,
				Opcode:                0x0,
				AuthoritativeAnswer:   false,
				Truncation:            false,
				RecursionDesired:      true,
				RecursionAvailable:    true,
				Reserved:              0x0,
				ResponseCode:          0x0,
				QuestionCount:         0x1,
				AnswerRecordCount:     0x6,
				AuthorityRecordCount:  0x0,
				AdditionalRecordCount: 0x0,
			},
		},
	}
	for i, c := range cases {
		got, err := toydns.ParseHeader(c.header)
		fmt.Printf("%#v \n", got)
		fmt.Printf("%#v \n", c.want)
		if err != nil {
			t.Fatalf("parse header error: %v", err)
		}
		if !reflect.DeepEqual(*got, c.want) {
			t.Fatalf("cases[%d] is error", i)
		}
	}
}
