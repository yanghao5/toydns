package toydns_test

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
	"toydns"
)

func TestQuestionToBytes(t *testing.T) {
	cases := []struct {
		question toydns.Question
		bytes    []byte
	}{
		{
			question: toydns.Question{
				QName:  "google.com",
				QType:  1,
				QClass: 1,
			},
			bytes: []byte{
				0x06, 'g', 'o', 'o', 'g', 'l', 'e',
				0x03, 'c', 'o', 'm',
				0x00,       // End of domain
				0x00, 0x01, // QTYPE = A
				0x00, 0x01, // QCLASS = IN
			},
		},
		{
			question: toydns.Question{
				QName:  "www.example.com",
				QType:  1,
				QClass: 1,
			},
			bytes: []byte{
				0x03, 'w', 'w', 'w',
				0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e',
				0x03, 'c', 'o', 'm',
				0x00,
				0x00, 0x01, // QTYPE = A
				0x00, 0x01, // QCLASS = IN
			},
		},
	}

	for i, c := range cases {
		if !bytes.Equal(c.bytes, c.question.ToBytes()) {
			t.Fatalf("cases[%d] not equal ", i)
		}
	}
}

func TestParseQuestion(t *testing.T) {
	cases := []struct {
		name     string
		buf      []byte
		question toydns.Question
		offset   int
	}{
		{
			name: "Google",
			buf: []byte{
				0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0,
				0x06, 'g', 'o', 'o', 'g', 'l', 'e',
				0x03, 'c', 'o', 'm',
				0x00,       // End of domain
				0x00, 0x01, // QTYPE = A
				0x00, 0x01, // QCLASS = IN
			},
			question: toydns.Question{
				QName:  "google.com",
				QType:  1,
				QClass: 1,
			},
			offset: 28,
		},
		{
			name: "Example",
			buf: []byte{
				0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0,
				0x03, 'w', 'w', 'w',
				0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e',
				0x03, 'c', 'o', 'm',
				0x00,
				0x00, 0x01, // QTYPE = A
				0x00, 0x01, // QCLASS = IN
			},
			question: toydns.Question{
				QName:  "www.example.com",
				QType:  1,
				QClass: 1,
			},
			offset: 33,
		},
	}
	for i, c := range cases {
		q, offset, err := toydns.ParseQuestion(c.buf)
		if err != nil {
			t.Fatalf("parse question error:%v \n", err)
		}
		if !reflect.DeepEqual(*q, c.question) {
			t.Fatalf("cases[%d] question %#v", i, q)
		}
		if offset != c.offset {
			t.Fatalf("cases[%d] offset %d", i, offset)
		}
	}
}

// go test -v -run TestEncodeDomain
func TestEncodeDomain(t *testing.T) {
	cases := []struct {
		name   string
		domain string
		encode []byte
	}{
		{
			name:   "Google",
			domain: "google.com",
			encode: []byte("\x06google\x03com\x00"),
		},
		{
			name:   "Facebook",
			domain: "facebook.com",
			encode: []byte("\x08facebook\x03com\x00"),
		},
		{
			name:   "dash domain",
			domain: "cdn-ch.yonomesh.com",
			encode: []byte("\x06cdn-ch\x08yonomesh\x03com\x00"),
		},
	}
	for i, c := range cases {
		if !bytes.Equal(toydns.EncodeDomain(c.domain), c.encode) {
			t.Fatalf("cases[%d] error", i)
		}
	}
}

// go test -v -run TestDecodeDomain
func TestDecodeDomain(t *testing.T) {
	cases := []struct {
		name   string
		domain []byte
		decode string
		move   int
	}{
		{
			name:   "Google",
			domain: []byte("\x06google\x03com\x00"),
			decode: "google.com",
			move:   12,
		},
		{
			name:   "Facebook",
			domain: []byte("\x08facebook\x03com\x00"),
			decode: "facebook.com",
			move:   14,
		},
		{
			name:   "dash domain",
			domain: []byte("\x06cdn-ch\x08yonomesh\x03com\x00"),
			decode: "cdn-ch.yonomesh.com",
			move:   21,
		},
	}
	for i, c := range cases {
		domain, move := toydns.DecodeDomain(c.domain)
		fmt.Println(domain)
		if domain != c.decode {
			t.Fatalf("cases[%d] domain error", i)
		}
		if move != c.move {
			t.Fatalf("cases[%d] move error", i)
		}
	}
}
