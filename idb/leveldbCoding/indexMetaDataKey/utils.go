package indexMetaDataKey

import (
	"idb-parser/idb/leveldbCoding"
	"idb-parser/idb/leveldbCoding/keyPrefix"
	"idb-parser/idb/leveldbCoding/varint"
)

type IndexMetaDataKey struct {
	ObjectStoreId int64
	IndexId       int64
	MetaDataType  MetaDataType
}

type MetaDataType uint8

const (
	Name MetaDataType = iota
	Unique
	KeyPath
	MultiEntry
)

func (k IndexMetaDataKey) Compare(other IndexMetaDataKey) int {
	if k.ObjectStoreId < 0 || k.IndexId < 0 {
		panic("k.ObjectStoreId < 0 || k.IndexId < 0")
	}
	if x := leveldbCoding.CompareInts(k.ObjectStoreId, other.ObjectStoreId); x != 0 {
		return x
	}
	if x := leveldbCoding.CompareInts(k.IndexId, other.IndexId); x != 0 {
		return x
	}
	return int(k.MetaDataType - other.MetaDataType)
}

func (k IndexMetaDataKey) Decode(slice *[]byte, result *IndexMetaDataKey) bool {
	var prefix keyPrefix.KeyPrefix
	if !(keyPrefix.KeyPrefix{}).Decode(slice, &prefix) {
		return false
	}
	if prefix.DatabaseId == 0 || prefix.ObjectStoreId != 0 || prefix.IndexId != 0 {
		panic("prefix.DatabaseId == 0 || prefix.ObjectStoreId != 0 && prefix.IndexId != 0")
	}
	var typeByte byte = 0
	if !leveldbCoding.DecodeByte(slice, &typeByte) {
		return false
	}
	if typeByte != leveldbCoding.KIndexMetaDataTypeByte {
		panic("typeByte != leveldbCoding.KIndexMetaDataTypeByte")
	}
	if !varint.DecodeVarInt(slice, &result.ObjectStoreId) {
		return false
	}
	if !varint.DecodeVarInt(slice, &result.IndexId) {
		return false
	}
	if !leveldbCoding.DecodeByte(slice, (*byte)(&result.MetaDataType)) {
		return false
	}
	return true
}

func (k IndexMetaDataKey) Encode(databaseId, objectStoreId, indexId int64, metaDataType MetaDataType) string {
	prefix := keyPrefix.KeyPrefix{
		DatabaseId:    databaseId,
		ObjectStoreId: 0,
		IndexId:       0,
	}
	ret := prefix.Encode()
	ret += string(leveldbCoding.KIndexMetaDataTypeByte)
	varint.EncodeVarInt(objectStoreId, &ret)
	varint.EncodeVarInt(indexId, &ret)
	leveldbCoding.EncodeByte(byte(metaDataType), &ret)
	return ret
}
