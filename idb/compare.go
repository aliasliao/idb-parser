package idb

import (
	"bytes"

	"idb-parser/idb/databaseFreeListKey"
	"idb-parser/idb/keyPrefix"
	"idb-parser/idb/leveldbCoding"
)

func Compare(a, b []byte, onlyCompareIndexKeys bool, ok *bool) int {
	sliceA := append([]byte{}, a...)
	sliceB := append([]byte{}, b...)
	prefixA := keyPrefix.KeyPrefix{}
	prefixB := keyPrefix.KeyPrefix{}

	okA := keyPrefix.KeyPrefix{}.Decode(&sliceA, &prefixA)
	okB := keyPrefix.KeyPrefix{}.Decode(&sliceB, &prefixB)
	if !okA || !okB {
		*ok = false
		return 0
	}

	*ok = true
	if x := prefixA.Compare(prefixB); x != 0 {
		return x
	}

	switch prefixA.Type() {
	case keyPrefix.GlobalMetadata:
		var typeByteA byte
		if !leveldbCoding.DecodeByte(&sliceA, &typeByteA) {
			*ok = false
			return 0
		}
		var typeByteB byte
		if !leveldbCoding.DecodeByte(&sliceB, &typeByteB) {
			*ok = false
			return 0
		}
		if x := int(typeByteA) - int(typeByteB); x != 0 {
			return x
		}
		if typeByteA < leveldbCoding.KMaxSimpleGlobalMetaDataTypeByte {
			return 0
		}
		if typeByteA == leveldbCoding.KScopesPrefixByte {
			return bytes.Compare(sliceA, sliceB)
		}
		if typeByteA == leveldbCoding.KDatabaseFreeListTypeByte {
			return leveldbCoding.CompareGeneric[databaseFreeListKey.DataBaseFreeListKey](sliceA, sliceB, onlyCompareIndexKeys, ok)
		}
		// TODO
	}

	return 1
}
