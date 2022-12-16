package compare

import (
	"bytes"

	"idb-parser/idb/leveldbCoding"
	"idb-parser/idb/leveldbCoding/blobEntryKey"
	"idb-parser/idb/leveldbCoding/databaseFreeListKey"
	"idb-parser/idb/leveldbCoding/databaseMetaDataKey"
	"idb-parser/idb/leveldbCoding/databaseNameKey"
	"idb-parser/idb/leveldbCoding/existsEntryKey"
	"idb-parser/idb/leveldbCoding/indexDataKey"
	"idb-parser/idb/leveldbCoding/indexFreeListKey"
	"idb-parser/idb/leveldbCoding/indexMetaDataKey"
	"idb-parser/idb/leveldbCoding/indexNamesKey"
	"idb-parser/idb/leveldbCoding/keyPrefix"
	"idb-parser/idb/leveldbCoding/objectStoreDataKey"
	"idb-parser/idb/leveldbCoding/objectStoreFreeListKey"
	"idb-parser/idb/leveldbCoding/objectStoreMetaDataKey"
	"idb-parser/idb/leveldbCoding/objectStoreNamesKey"
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

func CompareInternal(a, b []byte, onlyCompareIndexKeys bool, ok *bool) int {
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
			return bytes.Compare(sliceA, sliceB) // TODO: verify
		}
		if typeByteA == leveldbCoding.KDatabaseFreeListTypeByte {
			return CompareGeneric[databaseFreeListKey.DataBaseFreeListKey](a, b, onlyCompareIndexKeys, ok)
		}
		if typeByteA == leveldbCoding.KDatabaseNameTypeByte {
			return CompareGeneric[databaseNameKey.DatabaseNameKey](a, b, false, ok)
		}
	case keyPrefix.DatabaseMetadata:
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

		if typeByteA < byte(databaseMetaDataKey.MaxSimpleMetadataType) {
			return 0
		}
		if typeByteA == leveldbCoding.KObjectStoreMetaDataTypeByte {
			return CompareGeneric[objectStoreMetaDataKey.ObjectStoreMetaDataKey](a, b, onlyCompareIndexKeys, ok)
		}
		if typeByteA == leveldbCoding.KIndexMetaDataTypeByte {
			return CompareGeneric[indexMetaDataKey.IndexMetaDataKey](a, b, false, ok)
		}
		if typeByteA == leveldbCoding.KObjectStoreFreeListTypeByte {
			return CompareGeneric[objectStoreFreeListKey.ObjectStoreFreeListKey](a, b, onlyCompareIndexKeys, ok)
		}
		if typeByteA == leveldbCoding.KIndexFreeListTypeByte {
			return CompareGeneric[indexFreeListKey.IndexFreeListKey](a, b, false, ok)
		}
		if typeByteA == leveldbCoding.KObjectStoreNamesTypeByte {
			return CompareGeneric[objectStoreNamesKey.ObjectStoreNamesKey](a, b, onlyCompareIndexKeys, ok)
		}
		if typeByteA == leveldbCoding.KIndexNamesKeyTypeByte {
			return CompareGeneric[indexNamesKey.IndexNamesKey](a, b, false, ok)
		}
	case keyPrefix.ObjectStoreData:
		if len(sliceA) == 0 || len(sliceB) == 0 {
			return leveldbCoding.CompareSizes(len(sliceA), len(sliceB))
		}
		return objectStoreDataKey.CompareSuffix(&sliceA, &sliceB, false, ok)
	case keyPrefix.ExistsEntry:
		if len(sliceA) == 0 || len(sliceB) == 0 {
			return leveldbCoding.CompareSizes(len(sliceA), len(sliceB))
		}
		return existsEntryKey.CompareSuffix(&sliceA, &sliceB, false, ok)
	case keyPrefix.BlobEntry:
		if len(sliceA) == 0 || len(sliceB) == 0 {
			return leveldbCoding.CompareSizes(len(sliceA), len(sliceB))
		}
		return blobEntryKey.CompareSuffix(&sliceA, &sliceB, false, ok)
	case keyPrefix.IndexData:
		if len(sliceA) == 0 || len(sliceB) == 0 {
			return leveldbCoding.CompareSizes(len(sliceA), len(sliceB))
		}
		return indexDataKey.CompareSuffix(&sliceA, &sliceB, false, ok)
	case keyPrefix.InvalidType:
	}

	*ok = false
	return 0
}

func Compare(sliceA, sliceB []byte, onlyCompareIndexKeys bool) int {
	var ok bool
	result := CompareInternal(sliceA, sliceB, onlyCompareIndexKeys, &ok)
	if !ok {
		return 0
	}
	return result
}

func CompareKeys(sliceA, sliceB []byte) int {
	return Compare(sliceA, sliceB, false)
}
