package coding

import (
	"encoding/binary"

	"idb-parser/idb/coding/varint"
)

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
	varint.EncodeVarInt(len(value), into)
	EncodeString(value, into)
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

func DecodeInt(slice *[]byte, value *int64) bool {
	sliceValue := *slice
	if len(sliceValue) == 0 {
		return false
	}
	var ret int64 = 0
	for i, c := range sliceValue {
		ret |= c << (i * 8)
	}
	*value = ret

	sliceValue = (sliceValue)[len(sliceValue):]
	*slice = sliceValue
	return true
}

func CompareInts(a, b int64) int {
	diff := a - b
	if diff < 0 {
		return -1
	}
	if diff > 0 {
		return 1
	}
	return 0
}
