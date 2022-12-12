package leveldbCoding

import (
	"bytes"
)

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

func Compare(a, b []byte, onlyCompareIndexKeys bool, ok *bool) int {
	sliceA := append([]byte{}, a...)
	sliceB := append([]byte{}, b...)
	prefixA := KeyPrefix{}
	prefixB := KeyPrefix{}

	okA := KeyPrefix{}.Decode(&sliceA, &prefixA)
	okB := KeyPrefix{}.Decode(&sliceB, &prefixB)
	if !okA || !okB {
		*ok = false
		return 0
	}

	*ok = true
	if x := prefixA.Compare(prefixB); x != 0 {
		return x
	}

	switch prefixA.Type() {
	case GlobalMetadata:
		var typeByteA byte
		if !DecodeByte(&sliceA, &typeByteA) {
			*ok = false
			return 0
		}
		var typeByteB byte
		if !DecodeByte(&sliceB, &typeByteB) {
			*ok = false
			return 0
		}
		if x := int(typeByteA) - int(typeByteB); x != 0 {
			return x
		}
		if typeByteA < KMaxSimpleGlobalMetaDataTypeByte {
			return 0
		}
		if typeByteA == KScopesPrefixByte {
			return bytes.Compare(sliceA, sliceB) // TODO: verify
		}
		if typeByteA == KDatabaseFreeListTypeByte {
			return CompareGeneric[DataBaseFreeListKey](a, b, onlyCompareIndexKeys, ok)
		}
		if typeByteA == KDatabaseNameTypeByte {
			return CompareGeneric[DataBaseNameKey](a, b, false, ok)
		}
	}

	return 1
}
