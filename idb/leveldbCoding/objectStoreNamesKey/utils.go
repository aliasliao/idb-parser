package objectStoreNamesKey

import (
	"idb-parser/idb/leveldbCoding"
	"idb-parser/idb/leveldbCoding/keyPrefix"
)

type ObjectStoreNamesKey struct {
	ObjectStoreName leveldbCoding.U16string
}

func (k ObjectStoreNamesKey) Compare(other ObjectStoreNamesKey) int {
	return leveldbCoding.CompareU16String(k.ObjectStoreName, other.ObjectStoreName)
}

func (k ObjectStoreNamesKey) Decode(slice *[]byte, result *ObjectStoreNamesKey) bool {
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
	if typeByte != leveldbCoding.KObjectStoreNamesTypeByte {
		panic("typeByte != leveldbCoding.KObjectStoreNamesTypeByte")
	}
	if !leveldbCoding.DecodeStringWithLength(slice, &result.ObjectStoreName) {
		return false
	}
	return true
}
