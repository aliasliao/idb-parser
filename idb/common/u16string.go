package common

import (
	"bytes"
	"encoding/binary"
)

type U16string []uint16

func CompareU16String(a, b U16string) int {
	bytesA := make([]byte, len(a)*2)
	for i, c := range a {
		binary.BigEndian.PutUint16(bytesA[i*2:], c)
	}

	bytesB := make([]byte, len(b)*2)
	for i, c := range b {
		binary.BigEndian.PutUint16(bytesB[i*2:], c)
	}

	return bytes.Compare(bytesA, bytesB)
}

func (s U16string) ToString() string {
	buf := make([]byte, len(s))
	for i, c := range s {
		buf[i] = byte(c)
	}
	return string(buf)
}
