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
	kMaxDatabaseIdSizeBits     = 3
	kMaxObjectStoreIdSizeBits  = 3
	kMaxIndexIdSizeBits        = 2
	kMaxDatabaseIdSizeBytes    = 1 << kMaxDatabaseIdSizeBits    // 8
	kMaxObjectStoreIdSizeBytes = 1 << kMaxObjectStoreIdSizeBits // 8
	kMaxIndexIdSizeBytes       = 1 << kMaxIndexIdSizeBits       // 4

	kMaxDatabaseIdBits    = kMaxDatabaseIdSizeBytes*8 - 1    // 63
	kMaxObjectStoreIdBits = kMaxObjectStoreIdSizeBytes*8 - 1 // 63
	kMaxIndexIdBits       = kMaxIndexIdSizeBytes*8 - 1       // 31

	kMaxDatabaseId    int64 = (1 << kMaxDatabaseIdBits) - 1    // max signed int64_t
	kMaxObjectStoreId int64 = (1 << kMaxObjectStoreIdBits) - 1 // max signed int64_t
	kMaxIndexId       int64 = (1 << kMaxIndexIdBits) - 1       // max signed int32_t

	kInvalidId int64 = -1
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
	if k.DatabaseId == kInvalidId || k.ObjectStoreId == kInvalidId || k.IndexId == kInvalidId {
		panic("k.DatabaseId == kInvalidId || k.ObjectStoreId == kInvalidId || k.IndexId == kInvalidId")
	}
	return EncodeInternal(k.DatabaseId, k.ObjectStoreId, k.IndexId)
}

func EncodeInternal(databaseId, objectStoreId, indexId int64) string {
	var databaseIdStr string
	var objectSoreIdStr string
	var indexIdStr string

	leveldbCoding.EncodeIntSafely(databaseId, kMaxDatabaseId, &databaseIdStr)
	leveldbCoding.EncodeIntSafely(objectStoreId, kMaxObjectStoreId, &objectSoreIdStr)
	leveldbCoding.EncodeIntSafely(indexId, kMaxIndexId, &indexIdStr)

	if len(databaseIdStr) > kMaxDatabaseIdSizeBytes || len(objectSoreIdStr) > kMaxObjectStoreIdSizeBytes || len(indexIdStr) > kMaxIndexIdSizeBytes {
		panic("len(databaseIdStr) > kMaxDatabaseIdSizeBytes || len(objectSoreIdStr) > kMaxObjectStoreIdSizeBytes || len(indexIdStr) > kMaxIndexIdSizeBytes")
	}

	firstByteNum := (len(databaseIdStr)-1)<<(kMaxObjectStoreIdSizeBits+kMaxIndexIdSizeBits) | (len(objectSoreIdStr)-1)<<kMaxIndexIdSizeBits | (len(indexIdStr) - 1)
	firstByte := byte(firstByteNum)

	ret := string(firstByte)
	ret += databaseIdStr
	ret += objectSoreIdStr
	ret += indexIdStr
	return ret
}

func IsValidDatabaseId(databaseId int64) bool {
	return databaseId > 0 && databaseId < kMaxDatabaseId
}
