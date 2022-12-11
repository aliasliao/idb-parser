package coding

import "encoding/binary"

type U16string = []uint16

func ASCIIToUTF16(s string) U16string {
	list := []byte(s)
	ret := make(U16string, len(list))
	for i, b := range list {
		ret[i] = uint16(b)
	}
	return ret
}

func EncodeStringWithLength(value U16string, into *string) {
	EncodeVarInt(len(value), into)
	EncodeString(value, into)
}

func EncodeVarInt(from int, into *string) {
	n := uint64(from)
	var buf []byte
	if n == 0 {
		buf = []byte{0}
	}
	for n > 0 {
		c := byte(n & 0x7f)
		n >>= 7
		if n > 0 {
			c |= 0x80
		}
		buf = append(buf, c)
	}
	*into += string(buf)
}

func EncodeString(from U16string, into *string) {
	if len(from) == 0 {
		return
	}
	buf := make([]byte, len(from)*2)
	for i, c := range from {
		binary.BigEndian.PutUint16(buf[i*2:], c)
	}
	*into += string(buf)
}
