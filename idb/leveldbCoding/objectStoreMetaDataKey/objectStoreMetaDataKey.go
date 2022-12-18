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

const (
	KKeyGeneratorInitialNumber int64 = 1
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

	if !leveldbCoding.DecodeByte(slice, (*byte)(&result.MetaDataType)) {
		return false
	}
	return true
}

func (k ObjectStoreMetaDataKey) Encode(databaseId int64, objectStoreId int64, metaDataType MetaDataType) string {
	prefix := keyPrefix.KeyPrefix{DatabaseId: databaseId}
	ret := prefix.Encode()
	ret += string(leveldbCoding.KObjectStoreMetaDataTypeByte)
	varint.EncodeVarInt(objectStoreId, &ret)
	ret += string(metaDataType)
	return ret
}

func (k ObjectStoreMetaDataKey) EncodeMaxKey(databaseId int64) string {
	return k.Encode(databaseId, keyPrefix.KMaxObjectStoreId, MetaDataType(leveldbCoding.KObjectMetaDataTypeMaximum))
}
