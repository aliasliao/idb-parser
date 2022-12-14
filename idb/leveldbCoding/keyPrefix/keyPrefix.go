package keyPrefix

import "idb-parser/idb/leveldbCoding"

type KeyPrefix struct {
	DatabaseId    int64
	ObjectStoreId int64
	IndexId       int64
}

type KeyPrefixType uint8

const (
	GlobalMetadata KeyPrefixType = iota
	DatabaseMetadata
	ObjectStoreData
	ExistsEntry
	IndexData
	InvalidType
	BlobEntry
)
const (
	kMinimumIndexId byte = 30
)

func (k KeyPrefix) Type() KeyPrefixType {
	if k.DatabaseId == 0 {
		return GlobalMetadata
	}
	if k.ObjectStoreId == 0 {
		return DatabaseMetadata
	}
	if k.IndexId == int64(leveldbCoding.KObjectStoreDataIndexId) {
		return ObjectStoreData
	}
	if k.IndexId == int64(leveldbCoding.KExistsEntryIndexId) {
		return ExistsEntry
	}
	if k.IndexId == int64(leveldbCoding.KBlobEntryIndexId) {
		return BlobEntry
	}
	if k.IndexId >= int64(kMinimumIndexId) {
		return IndexData
	}
	return InvalidType
}

func (k KeyPrefix) Compare(other KeyPrefix) int {
	if k.DatabaseId != other.DatabaseId {
		return leveldbCoding.CompareInts(k.DatabaseId, other.DatabaseId)
	}
	if k.ObjectStoreId != other.ObjectStoreId {
		return leveldbCoding.CompareInts(k.ObjectStoreId, other.ObjectStoreId)
	}
	if k.IndexId != other.IndexId {
		return leveldbCoding.CompareInts(k.IndexId, other.IndexId)
	}
	return 0
}

func (k KeyPrefix) Decode(slice *[]byte, result *KeyPrefix) bool {
	sliceValue := *slice
	if len(sliceValue) == 0 {
		return false
	}
	firstByte := sliceValue[0]
	sliceValue = sliceValue[1:]

	databaseIdBytes := int(((firstByte >> 5) & 0x7) + 1)
	objectStoreIdBytes := int(((firstByte >> 2) & 0x7) + 1)
	indexIdBytes := int(firstByte&0x3 + 1)

	if databaseIdBytes+objectStoreIdBytes+indexIdBytes > len(sliceValue) {
		return false
	}

	{
		tmp := sliceValue[0:databaseIdBytes]
		if !leveldbCoding.DecodeInt(&tmp, &(result.DatabaseId)) {
			return false
		}
	}
	sliceValue = sliceValue[databaseIdBytes:]
	{
		tmp := sliceValue[0:objectStoreIdBytes]
		if !leveldbCoding.DecodeInt(&tmp, &(result.ObjectStoreId)) {
			return false
		}
	}
	sliceValue = sliceValue[objectStoreIdBytes:]
	{
		tmp := sliceValue[0:indexIdBytes]
		if !leveldbCoding.DecodeInt(&tmp, &(result.IndexId)) {
			return false
		}
	}
	sliceValue = sliceValue[indexIdBytes:]
	*slice = sliceValue
	return true
}
