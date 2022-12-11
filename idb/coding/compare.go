package coding

import (
	"bytes"

	"idb-parser/idb/databaseFreeListKey"
	"idb-parser/idb/keyPrefix"
)

func Compare(a, b []byte, onlyCompareIndexKeys bool) int {
	sliceA := append([]byte{}, a...)
	sliceB := append([]byte{}, b...)
	prefixA := keyPrefix.KeyPrefix{}
	prefixB := keyPrefix.KeyPrefix{}

	okA := keyPrefix.KeyPrefix{}.Decode(&sliceA, &prefixA)
	okB := keyPrefix.KeyPrefix{}.Decode(&sliceB, &prefixB)
	if !okA || !okB {
		return 0
	}

	if x := prefixA.Compare(prefixB); x != 0 {
		return x
	}

	switch prefixA.Type() {
	case keyPrefix.GlobalMetadata:
		var typeByteA byte
		if !DecodeByte(&sliceA, &typeByteA) {
			return 0
		}
		var typeByteB byte
		if !DecodeByte(&sliceB, &typeByteB) {
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
			return CompareGeneric[databaseFreeListKey.DataBaseFreeListKey](sliceA, sliceB, onlyCompareIndexKeys)
		}
	}

	return 1
}
