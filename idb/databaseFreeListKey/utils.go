package databaseFreeListKey

import (
	"idb-parser/idb/coding"
	"idb-parser/idb/coding/varint"
	"idb-parser/idb/keyPrefix"
)

type DataBaseFreeListKey struct {
	databaseId int64
}

func (k DataBaseFreeListKey) Compare(other DataBaseFreeListKey) int {
	return coding.CompareInts(k.databaseId, other.databaseId)
}

func (DataBaseFreeListKey) Decode(slice *[]byte, k *DataBaseFreeListKey) bool {
	var prefix keyPrefix.KeyPrefix
	if !(keyPrefix.KeyPrefix{}.Decode(slice, &prefix)) {
		return false
	}
	var typeByte byte
	if !coding.DecodeByte(slice, &typeByte) {
		return false
	}
	if !varint.DecodeVarInt(slice, &k.databaseId) {
		return false
	}
	return true
}
