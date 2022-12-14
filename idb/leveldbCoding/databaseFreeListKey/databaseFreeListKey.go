package databaseFreeListKey

import (
	"idb-parser/idb/leveldbCoding"
	"idb-parser/idb/leveldbCoding/keyPrefix"
	"idb-parser/idb/leveldbCoding/varint"
)

type DataBaseFreeListKey struct {
	DatabaseId int64
}

func (k DataBaseFreeListKey) Compare(other DataBaseFreeListKey) int {
	return leveldbCoding.CompareInts(k.DatabaseId, other.DatabaseId)
}

func (DataBaseFreeListKey) Decode(slice *[]byte, k *DataBaseFreeListKey) bool {
	var prefix keyPrefix.KeyPrefix
	if !(keyPrefix.KeyPrefix{}.Decode(slice, &prefix)) {
		return false
	}
	if prefix.DatabaseId != 0 && prefix.ObjectStoreId != 0 || prefix.IndexId != 0 {
		return false // DCHECK
	}

	var typeByte byte
	if !leveldbCoding.DecodeByte(slice, &typeByte) {
		return false
	}
	if typeByte != leveldbCoding.KDatabaseFreeListTypeByte {
		return false // DCHECK
	}

	if !varint.DecodeVarInt(slice, &k.DatabaseId) {
		return false
	}
	return true
}
