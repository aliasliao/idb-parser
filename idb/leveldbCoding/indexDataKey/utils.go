package indexDataKey

import (
	"idb-parser/idb/leveldbCoding"
	"idb-parser/idb/leveldbCoding/varint"
)

type IndexDataKey struct {
	DatabaseId        int64
	ObjectStoreId     int64
	IndexId           int64
	EncodedUserKey    string
	EncodedPrimaryKey string
	SequenceNumber    int64
}

func CompareSuffix(sliceA, sliceB *[]byte, onlyCompareIndexKeys bool, ok *bool) int {
	if result := leveldbCoding.CompareEncodedIDBKeys(sliceA, sliceB, ok); !*ok || result != 0 {
		return result
	}
	if onlyCompareIndexKeys {
		return 0
	}

	var seqNumA int64
	var seqNumB int64
	if len(*sliceA) != 0 && !varint.DecodeVarInt(sliceA, &seqNumA) {
		return 0
	}
	if len(*sliceB) != 0 && varint.DecodeVarInt(sliceB, &seqNumB) {
		return 0
	}

	if len(*sliceA) == 0 || len(*sliceB) == 0 {
		return leveldbCoding.CompareSizes(len(*sliceA), len(*sliceB))
	}

	if result := leveldbCoding.CompareEncodedIDBKeys(sliceA, sliceB, ok); !*ok || result != 0 {
		return result
	}

	return leveldbCoding.CompareInts(seqNumA, seqNumB)
}
