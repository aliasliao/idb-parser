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
	KMaxDatabaseIdSizeBits     = 3
	KMaxObjectStoreIdSizeBits  = 3
	KMaxIndexIdSizeBits        = 2
	KMaxDatabaseIdSizeBytes    = 1 << KMaxDatabaseIdSizeBits    // 8
	KMaxObjectStoreIdSizeBytes = 1 << KMaxObjectStoreIdSizeBits // 8
	KMaxIndexIdSizeBytes       = 1 << KMaxIndexIdSizeBits       // 4

	KMaxDatabaseIdBits    = KMaxDatabaseIdSizeBytes*8 - 1    // 63
	KMaxObjectStoreIdBits = KMaxObjectStoreIdSizeBytes*8 - 1 // 63
	KMaxIndexIdBits       = KMaxIndexIdSizeBytes*8 - 1       // 31

	KMaxDatabaseId    int64 = (1 << KMaxDatabaseIdBits) - 1    // max signed int64_t
	KMaxObjectStoreId int64 = (1 << KMaxObjectStoreIdBits) - 1 // max signed int64_t
	KMaxIndexId       int64 = (1 << KMaxIndexIdBits) - 1       // max signed int32_t

	KInvalidId int64 = -1
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
	if k.IndexId >= int64(leveldbCoding.KMinimumIndexId) {
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

func (k KeyPrefix) Encode() string {
	if k.DatabaseId == KInvalidId || k.ObjectStoreId == KInvalidId || k.IndexId == KInvalidId {
		panic("k.DatabaseId == KInvalidId || k.ObjectStoreId == KInvalidId || k.IndexId == KInvalidId")
	}
	return EncodeInternal(k.DatabaseId, k.ObjectStoreId, k.IndexId)
}

func EncodeInternal(databaseId, objectStoreId, indexId int64) string {
	var databaseIdStr string
	var objectSoreIdStr string
	var indexIdStr string

	leveldbCoding.EncodeIntSafely(databaseId, KMaxDatabaseId, &databaseIdStr)
	leveldbCoding.EncodeIntSafely(objectStoreId, KMaxObjectStoreId, &objectSoreIdStr)
	leveldbCoding.EncodeIntSafely(indexId, KMaxIndexId, &indexIdStr)

	if len(databaseIdStr) > KMaxDatabaseIdSizeBytes || len(objectSoreIdStr) > KMaxObjectStoreIdSizeBytes || len(indexIdStr) > KMaxIndexIdSizeBytes {
		panic("len(databaseIdStr) > KMaxDatabaseIdSizeBytes || len(objectSoreIdStr) > KMaxObjectStoreIdSizeBytes || len(indexIdStr) > KMaxIndexIdSizeBytes")
	}

	firstByteNum := (len(databaseIdStr)-1)<<(KMaxObjectStoreIdSizeBits+KMaxIndexIdSizeBits) | (len(objectSoreIdStr)-1)<<KMaxIndexIdSizeBits | (len(indexIdStr) - 1)
	firstByte := byte(firstByteNum)

	ret := string(firstByte)
	ret += databaseIdStr
	ret += objectSoreIdStr
	ret += indexIdStr
	return ret
}

func IsValidDatabaseId(databaseId int64) bool {
	return databaseId > 0 && databaseId < KMaxDatabaseId
}

func IsValidObjectStoreId(objectStoreId int64) bool {
	return objectStoreId > 0 && objectStoreId < KMaxObjectStoreId
}

func ValidIds(databaseId, objectStoreId int64) bool {
	return IsValidDatabaseId(databaseId) && IsValidObjectStoreId(objectStoreId)
}
