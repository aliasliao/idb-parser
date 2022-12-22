package indexNamesKey

import (
	"idb-parser/idb/common"
	"idb-parser/idb/leveldbCoding"
	"idb-parser/idb/leveldbCoding/keyPrefix"
	"idb-parser/idb/leveldbCoding/varint"
)

type IndexNamesKey struct {
	ObjectStoreId int64
	IndexName     common.U16string
}

func (k IndexNamesKey) Compare(other IndexNamesKey) int {
	if k.ObjectStoreId <= 0 {
		panic("k.ObjectStoreId <= 0")
	}
	if x := leveldbCoding.CompareInts(k.ObjectStoreId, other.ObjectStoreId); x != 0 {
		return x
	}
	return common.CompareU16String(k.IndexName, other.IndexName)
}

func (k IndexNamesKey) Decode(slice *[]byte, result *IndexNamesKey) bool {
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
	if typeByte != leveldbCoding.KIndexNamesKeyTypeByte {
		panic("typeByte != leveldbCoding.KIndexNamesKeyTypeByte")
	}
	if !varint.DecodeVarInt(slice, &result.ObjectStoreId) {
		return false
	}
	if !leveldbCoding.DecodeStringWithLength(slice, &result.IndexName) {
		return false
	}
	return true
}
