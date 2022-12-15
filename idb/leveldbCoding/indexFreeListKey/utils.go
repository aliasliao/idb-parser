package indexFreeListKey

import (
	"idb-parser/idb/leveldbCoding"
	"idb-parser/idb/leveldbCoding/keyPrefix"
	"idb-parser/idb/leveldbCoding/varint"
)

type IndexFreeListKey struct {
	ObjectStoreId int64
	IndexId       int64
}

func (k IndexFreeListKey) Compare(other IndexFreeListKey) int {
	if k.ObjectStoreId <= 0 || k.IndexId <= 0 {
		panic("k.ObjectStoreId <= 0 || k.IndexId <= 0")
	}
	if x := leveldbCoding.CompareInts(k.ObjectStoreId, other.ObjectStoreId); x != 0 {
		return x
	}
	return leveldbCoding.CompareInts(k.IndexId, other.IndexId)
}

func (k IndexFreeListKey) Decode(slice *[]byte, result *IndexFreeListKey) bool {
	var prefix keyPrefix.KeyPrefix
	if !(keyPrefix.KeyPrefix{}).Decode(slice, &prefix) {
		return false
	}
	if prefix.DatabaseId == 0 || prefix.ObjectStoreId != 0 && prefix.IndexId != 0 {
		panic("prefix.DatabaseId == 0 || prefix.ObjectStoreId != 0 && prefix.IndexId != 0")
	}
	var typeByte byte = 0
	if !leveldbCoding.DecodeByte(slice, &typeByte) {
		return false
	}
	if typeByte != leveldbCoding.KIndexFreeListTypeByte {
		panic("typeByte != leveldbCoding.KIndexFreeListTypeByte")
	}
	if !varint.DecodeVarInt(slice, &result.ObjectStoreId) {
		return false
	}
	if !varint.DecodeVarInt(slice, &result.IndexId) {
		return false
	}
	return true
}
