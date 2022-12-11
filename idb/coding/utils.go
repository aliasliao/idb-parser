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
	varint.EncodeVarInt(int64(len(value)), into)
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
		ret |= int64(c) << (i * 8)
	}
	*value = ret

	*slice = sliceValue[len(sliceValue):]
	return true
}

func DecodeByte(slice *[]byte, value *byte) bool {
	sliceValue := *slice
	if len(sliceValue) == 0 {
		return false
	}
	*value = sliceValue[0]
	*slice = sliceValue[1:]
	return true
}

func DecodeBool(slice *[]byte, value *bool) bool {
	sliceValue := *slice
	if len(sliceValue) == 0 {
		return false
	}
	*value = sliceValue[0] != 0
	*slice = sliceValue[1:]
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

type KeyType[T interface{}] interface {
	Compare(other T) int
	Decode(*[]byte, *T) bool
}

func CompareGeneric[T KeyType[T]](a, b []byte, onlyCompareIndexKeys bool, ok *bool) int {
	var tmp T

	var keyA T
	sliceA := append([]byte{}, a...)
	if !tmp.Decode(&sliceA, &keyA) {
		*ok = false
		return 0
	}

	var keyB T
	sliceB := append([]byte{}, b...)
	if !tmp.Decode(&sliceB, &keyB) {
		*ok = false
		return 0
	}

	*ok = true
	return keyA.Compare(keyB)
}
