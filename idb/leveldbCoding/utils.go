package leveldbCoding

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"

	"idb-parser/idb/leveldbCoding/mojom"
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
	if !varint.DecodeVarInt(&sliceValue, &strLen) || strLen < 0 {
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
		ret[i] = binary.BigEndian.Uint16(sliceValue[i*2:])
	}
	*value = ret
	*slice = sliceValue[len(sliceValue):]
	return true
}

func EncodeIntSafely(value, max int64, into *string) {
	if value > max {
		panic("value > max")
	}
	EncodeInt(value, into)
}

func EncodeInt(value int64, into *string) {
	if value < 0 {
		panic("value < 0")
	}
	n := uint64(value)
	if n == 0 {
		*into = string(0)
		return
	}

	var ret string
	for n > 0 {
		c := uint8(n)
		ret += string(c)
		n >>= 8
	}
	*into = ret
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

func DecodeDouble(slice *[]byte, value *float64) bool {
	sliceValue := *slice
	if len(sliceValue) < 8 {
		return false
	}
	bits := binary.BigEndian.Uint64(sliceValue)
	*value = math.Float64frombits(bits)
	*slice = sliceValue[8:]
	return true
}

func DecodeBinary(slice *[]byte, value *[]byte) bool {
	sliceValue := *slice
	var binLen int64 = 0
	if !varint.DecodeVarInt(&sliceValue, &binLen) || binLen < 0 { // fixed
		return false
	}
	if len(sliceValue) < int(binLen) {
		return false
	}
	*value = sliceValue[0:binLen]
	*slice = sliceValue[binLen:]
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

func CompareTypes(a, b mojom.IDBKeyType) int {
	return int(b - a)
}

func CompareEncodedBinary(sliceA, sliceB *[]byte, ok *bool) int {
	var binA []byte
	var binB []byte
	if !DecodeBinary(sliceA, &binA) || !DecodeBinary(sliceB, &binB) {
		*ok = false
		return 0
	}
	*ok = true
	return bytes.Compare(binA, binB)
}

func CompareEncodedStringWithLength(sliceA, sliceB *[]byte, ok *bool) int {
	var strA U16string
	var strB U16string
	if !DecodeStringWithLength(sliceA, &strA) || !DecodeStringWithLength(sliceB, &strB) {
		*ok = false
		return 0
	}
	*ok = true
	return CompareU16String(strA, strB)
}

func KeyTypeByteToKeyType(t byte) mojom.IDBKeyType {
	switch t {
	case KIndexedDBKeyNullTypeByte:
		return mojom.Invalid
	case KIndexedDBKeyArrayTypeByte:
		return mojom.Array
	case KIndexedDBKeyBinaryTypeByte:
		return mojom.Binary
	case KIndexedDBKeyStringTypeByte:
		return mojom.String
	case KIndexedDBKeyDateTypeByte:
		return mojom.Date
	case KIndexedDBKeyNumberTypeByte:
		return mojom.Number
	case KIndexedDBKeyMinKeyTypeByte:
		return mojom.Min
	}

	panic(fmt.Sprintf("Get invalid type %v", t))
}

func CompareEncodedIDBKeys(sliceA, sliceB *[]byte, ok *bool) int {
	*ok = true
	var typeByteA byte = 0
	if !DecodeByte(sliceA, &typeByteA) {
		*ok = false
		return 0
	}
	var typeByteB byte = 0
	if !DecodeByte(sliceB, &typeByteB) {
		*ok = false
		return 0
	}

	if x := CompareTypes(KeyTypeByteToKeyType(typeByteA), KeyTypeByteToKeyType(typeByteB)); x != 0 {
		return x
	}

	switch typeByteA {
	case KIndexedDBKeyNullTypeByte:
	case KIndexedDBKeyMinKeyTypeByte:
		return 0
	case KIndexedDBKeyArrayTypeByte:
		{
			var lenA int64 = 0
			var lenB int64 = 0
			if !varint.DecodeVarInt(sliceA, &lenA) || !varint.DecodeVarInt(sliceB, &lenB) {
				*ok = false
				return 0
			}
			for i := int64(0); i < lenA && i < lenB; i += 1 {
				if result := CompareEncodedIDBKeys(sliceA, sliceB, ok); !*ok || result != 0 {
					return result
				}
			}
			return int(lenA - lenB)
		}
	case KIndexedDBKeyBinaryTypeByte:
		return CompareEncodedBinary(sliceA, sliceB, ok)
	case KIndexedDBKeyStringTypeByte:
		return CompareEncodedStringWithLength(sliceA, sliceB, ok)
	case KIndexedDBKeyDateTypeByte:
	case KIndexedDBKeyNumberTypeByte:
		{
			var numA float64 = 0
			var numB float64 = 0
			if !DecodeDouble(sliceA, &numA) || !DecodeDouble(sliceB, &numB) {
				*ok = false
				return 0
			}
			if numA < numB {
				return -1
			}
			if numA > numB {
				return 1
			}
			return 0
		}
	}

	panic(fmt.Sprintf("Get invalid type %v", typeByteA))
}
