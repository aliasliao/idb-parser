package leveldbCoding

import (
	"bytes"
)

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
			return bytes.Compare(sliceA, sliceB)
		}
		if typeByteA == KDatabaseFreeListTypeByte {
			return CompareGeneric[DataBaseFreeListKey](sliceA, sliceB, onlyCompareIndexKeys, ok)
		}
		// TODO
	}

	return 1
}
