package keyPrefix

import "idb-parser/idb/coding"

type KeyPrefix struct {
	databaseId    int64
	objectStoreId int64
	indexId       int64
}

type Type uint8

const (
	GlobalMetadata   Type = 0
	DatabaseMetadata      = 1
	ObjectStoreData       = 2
	ExistsEntry           = 3
	IndexData             = 4
	InvalidType           = 5
	BlobEntry             = 6
)
const (
	kMinimumIndexId byte = 30
)

func (kp KeyPrefix) Type() Type {
	if kp.databaseId == 0 {
		return GlobalMetadata
	}
	if kp.objectStoreId == 0 {
		return DatabaseMetadata
	}
	if kp.indexId == int64(coding.KObjectStoreDataIndexId) {
		return ObjectStoreData
	}
	if kp.indexId == int64(coding.KExistsEntryIndexId) {
		return ExistsEntry
	}
	if kp.indexId == int64(coding.KBlobEntryIndexId) {
		return BlobEntry
	}
	if kp.indexId >= int64(kMinimumIndexId) {
		return IndexData
	}
	return InvalidType
}

func (kp KeyPrefix) Compare(other KeyPrefix) int {
	if kp.databaseId != other.databaseId {
		return coding.CompareInts(kp.databaseId, other.databaseId)
	}
	if kp.objectStoreId != other.objectStoreId {
		return coding.CompareInts(kp.objectStoreId, other.objectStoreId)
	}
	if kp.indexId != other.indexId {
		return coding.CompareInts(kp.indexId, other.indexId)
	}
	return 0
}

func (KeyPrefix) Decode(slice *[]byte, result *KeyPrefix) bool {
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
		if !coding.DecodeInt(&tmp, &(result.databaseId)) {
			return false
		}
	}
	sliceValue = sliceValue[databaseIdBytes:]
	{
		tmp := sliceValue[0:objectStoreIdBytes]
		if !coding.DecodeInt(&tmp, &(result.objectStoreId)) {
			return false
		}
	}
	sliceValue = sliceValue[objectStoreIdBytes:]
	{
		tmp := sliceValue[0:indexIdBytes]
		if !coding.DecodeInt(&tmp, &(result.indexId)) {
			return false
		}
	}
	sliceValue = sliceValue[indexIdBytes:]
	*slice = sliceValue
	return true
}
