package leveldbCoding

import (
	"bytes"

	"idb-parser/idb/leveldbCoding/databaseFreeListKey"
	"idb-parser/idb/leveldbCoding/databaseMetaDataKey"
	"idb-parser/idb/leveldbCoding/databaseNameKey"
	"idb-parser/idb/leveldbCoding/keyPrefix"
	"idb-parser/idb/leveldbCoding/objectStoreMetaDataKey"
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
			return CompareGeneric[databaseFreeListKey.DataBaseFreeListKey](a, b, onlyCompareIndexKeys, ok)
		}
		if typeByteA == KDatabaseNameTypeByte {
			return CompareGeneric[databaseNameKey.DatabaseNameKey](a, b, false, ok)
		}
	case keyPrefix.DatabaseMetadata:
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

		if typeByteA < byte(databaseMetaDataKey.MaxSimpleMetadataType) {
			return 0
		}
		if typeByteA == KObjectStoreMetaDataTypeByte {
			return CompareGeneric[objectStoreMetaDataKey.ObjectStoreMetaDataKey](a, b, onlyCompareIndexKeys, ok)
		}
		if typeByteA == KIndexMetaDataTypeByte {
			return CompareGeneric[IndexMetaDataKey](a, b, false, ok)
		}
		if typeByteA == KObjectStoreFreeListTypeByte {
			return CompareGeneric[ObjectStoreFreeListKey](a, b, onlyCompareIndexKeys, ok)
		}
		if typeByteA == KIndexFreeListTypeByte {
			return CompareGeneric[IndexFreeListKey](a, b, false, ok)
		}
		if typeByteA == KObjectStoreNamesTypeByte {
			return CompareGeneric[ObjectStoreNamesKey](a, b, onlyCompareIndexKeys, ok)
		}
		if typeByteA == KIndexNamesKeyTypeByte {
			return CompareGeneric[IndexNamesKey](a, b, false, ok)
		}
	case keyPrefix.ObjectStoreData:
		if len(sliceA) == 0 || len(sliceB) == 0 {
			return CompareSizes(len(sliceA), len(sliceB))
		}
		return CompareSuffix[ObjectStoreDataKey](&sliceA, &sliceB, false, ok)
	case keyPrefix.ExistsEntry:
		if len(sliceA) == 0 || len(sliceB) == 0 {
			return CompareSizes(len(sliceA), len(sliceB))
		}
		return CompareSuffix[ExistsExtryKey](&sliceA, &sliceB, false, ok)
	case keyPrefix.BlobEntry:
		if len(sliceA) == 0 || len(sliceB) == 0 {
			return CompareSizes(len(sliceA), len(sliceB))
		}
		return CompareSuffix[BlobEntryKey](&sliceA, &sliceB, false, ok)
	case keyPrefix.IndexData:
		if len(sliceA) == 0 || len(sliceB) == 0 {
			return CompareSizes(len(sliceA), len(sliceB))
		}
		return CompareSuffix[IndexDataKey](&sliceA, &sliceB, false, ok)
	case keyPrefix.InvalidType:
	}

	*ok = false
	return 0
}
