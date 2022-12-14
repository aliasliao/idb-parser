package objectStoreFreeListKey

import (
	"idb-parser/idb/leveldbCoding"
	"idb-parser/idb/leveldbCoding/keyPrefix"
	"idb-parser/idb/leveldbCoding/varint"
)

type ObjectStoreFreeListKey struct {
	ObjectStoreId int64
}

func (k ObjectStoreFreeListKey) Compare(other ObjectStoreFreeListKey) int {
	if k.ObjectStoreId != 0 {
		panic("k.ObjectStoreId != 0")
	}
	return leveldbCoding.CompareInts(k.ObjectStoreId, other.ObjectStoreId)
}

func (k ObjectStoreFreeListKey) Decode(slice *[]byte, result *ObjectStoreFreeListKey) bool {
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
	if typeByte != leveldbCoding.KObjectStoreFreeListTypeByte {
		panic("typeByte != leveldbCoding.KObjectStoreFreeListTypeByte")
	}
	if !varint.DecodeVarInt(slice, &result.ObjectStoreId) {
		return false
	}
	return true
}
