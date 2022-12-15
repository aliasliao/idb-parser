package leveldbCoding

import (
	"encoding/binary"

	"idb-parser/idb/leveldbCoding/varint"
)

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

func DecodeStringWithLength(slice *[]byte, value *U16string) bool {
	sliceValue := *slice
	if len(sliceValue) == 0 {
		return false
	}
	var strLen int64 = 0
	if !varint.DecodeVarInt(&sliceValue, &strLen) || strLen == 0 {
		return false
	}
	bytesLen := int(strLen) * 2
	if len(sliceValue) < bytesLen {
		return false
	}
	buf := sliceValue[0:bytesLen]
	if !DecodeString(&buf, value) {
		return false
	}
	*slice = sliceValue[bytesLen:]
	return true
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

func DecodeString(slice *[]byte, value *U16string) bool {
	sliceValue := *slice
	bytesLen := len(sliceValue)
	if bytesLen == 0 {
		*value = U16string{}
		return true
	}
	if bytesLen%2 != 0 {
		return false // DCHECK
	}
	strLen := bytesLen / 2
	ret := make(U16string, strLen)
	for i := 0; i < strLen; i++ {
		ret[i] = binary.BigEndian.Uint16(sliceValue[i:])
	}
	*value = ret
	*slice = sliceValue[len(sliceValue):]
	return true
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

func CompareSizes(a, b int) int {
	if a > b {
		return 1
	}
	if a < b {
		return -1
	}
	return 0
}
