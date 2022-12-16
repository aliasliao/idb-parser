package blobEntryKey

import "idb-parser/idb/leveldbCoding"

var KSpecialIndexNumber int64

type BlobEntryKey struct {
	DatabaseId     int64
	ObjectStoreId  int64
	EncodedUserKey string
}

func CompareSuffix(sliceA, sliceB *[]byte, onlyCompareIndexKeys bool, ok *bool) int {
	return leveldbCoding.CompareEncodedIDBKeys(sliceA, sliceB, ok)
}
