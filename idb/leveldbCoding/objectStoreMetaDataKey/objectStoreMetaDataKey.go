package objectStoreMetaDataKey

import (
	"idb-parser/idb/leveldbCoding"
	"idb-parser/idb/leveldbCoding/keyPrefix"
	"idb-parser/idb/leveldbCoding/varint"
)

type ObjectStoreMetaDataKey struct {
	ObjectStoreId int64
	MetaDataType  MetaDataType
}

type MetaDataType uint8

const (
	Name MetaDataType = iota
	KeyPath
	AutoIncrement
	Evictable
	LastVersion
	MaxIndexId
	HasKeyPath
	KeyGeneratorCurrentNumber
)

func (k ObjectStoreMetaDataKey) Compare(other ObjectStoreMetaDataKey) int {
	if x := leveldbCoding.CompareInts(k.ObjectStoreId, other.ObjectStoreId); x != 0 {
		return x
	}
	return int(k.MetaDataType - other.MetaDataType)
}

func (k ObjectStoreMetaDataKey) Decode(slice *[]byte, result *ObjectStoreMetaDataKey) bool {
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
	if typeByte != leveldbCoding.KObjectStoreMetaDataTypeByte {
		panic("typeByte != leveldbCoding.KObjectStoreMetaDataTypeByte")
	}
	if !varint.DecodeVarInt(slice, &result.ObjectStoreId) {
		return false
	}
	if result.ObjectStoreId == 0 {
		panic("result.ObjectStoreId == 0")
	}

	var tmpByte byte
	if !leveldbCoding.DecodeByte(slice, &tmpByte) {
		return false
	}
	result.MetaDataType = MetaDataType(tmpByte)
	return true
}
