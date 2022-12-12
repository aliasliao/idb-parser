package leveldbCoding

import (
	"idb-parser/idb/leveldbCoding/varint"
)

type DataBaseFreeListKey struct {
	databaseId int64
}

func (k DataBaseFreeListKey) Compare(other DataBaseFreeListKey) int {
	return CompareInts(k.databaseId, other.databaseId)
}

func (DataBaseFreeListKey) Decode(slice *[]byte, k *DataBaseFreeListKey) bool {
	var prefix KeyPrefix
	if !(KeyPrefix{}.Decode(slice, &prefix)) {
		return false
	}
	var typeByte byte
	if !DecodeByte(slice, &typeByte) {
		return false
	}
	if !varint.DecodeVarInt(slice, &k.databaseId) {
		return false
	}
	return true
}