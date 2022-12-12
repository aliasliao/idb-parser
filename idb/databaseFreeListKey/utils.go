package databaseFreeListKey

import (
	"idb-parser/idb/keyPrefix"
	"idb-parser/idb/leveldbCoding"
	"idb-parser/idb/leveldbCoding/varint"
)

type DataBaseFreeListKey struct {
	databaseId int64
}

func (k DataBaseFreeListKey) Compare(other DataBaseFreeListKey) int {
	return leveldbCoding.CompareInts(k.databaseId, other.databaseId)
}

func (DataBaseFreeListKey) Decode(slice *[]byte, k *DataBaseFreeListKey) bool {
	var prefix keyPrefix.KeyPrefix
	if !(keyPrefix.KeyPrefix{}.Decode(slice, &prefix)) {
		return false
	}
	var typeByte byte
	if !leveldbCoding.DecodeByte(slice, &typeByte) {
		return false
	}
	if !varint.DecodeVarInt(slice, &k.databaseId) {
		return false
	}
	return true
}
