package existsEntryKey

import "idb-parser/idb/leveldbCoding"

var KSpecialIndexNumber int64

type ExistsEntryKey struct {
	EncodedUserKey string
}

func CompareSuffix(sliceA, sliceB *[]byte, onlyCompareIndexKeys bool, ok *bool) int {
	return leveldbCoding.CompareEncodedIDBKeys(sliceA, sliceB, ok)
}
