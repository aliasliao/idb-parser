package objectStoreDataKey

import "idb-parser/idb/leveldbCoding"

type ObjectStoreDataKey struct {
	EncodedUserKey string
}

func CompareSuffix(sliceA, sliceB *[]byte, onlyCompareIndexKeys bool, ok *bool) int {
	return leveldbCoding.CompareEncodedIDBKeys(sliceA, sliceB, ok)
}
