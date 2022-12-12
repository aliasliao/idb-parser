package leveldbCoding

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

func (k KeyPrefix) Type() Type {
	if k.databaseId == 0 {
		return GlobalMetadata
	}
	if k.objectStoreId == 0 {
		return DatabaseMetadata
	}
	if k.indexId == int64(KObjectStoreDataIndexId) {
		return ObjectStoreData
	}
	if k.indexId == int64(KExistsEntryIndexId) {
		return ExistsEntry
	}
	if k.indexId == int64(KBlobEntryIndexId) {
		return BlobEntry
	}
	if k.indexId >= int64(kMinimumIndexId) {
		return IndexData
	}
	return InvalidType
}

func (k KeyPrefix) Compare(other KeyPrefix) int {
	if k.databaseId != other.databaseId {
		return CompareInts(k.databaseId, other.databaseId)
	}
	if k.objectStoreId != other.objectStoreId {
		return CompareInts(k.objectStoreId, other.objectStoreId)
	}
	if k.indexId != other.indexId {
		return CompareInts(k.indexId, other.indexId)
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
		if !DecodeInt(&tmp, &(result.databaseId)) {
			return false
		}
	}
	sliceValue = sliceValue[databaseIdBytes:]
	{
		tmp := sliceValue[0:objectStoreIdBytes]
		if !DecodeInt(&tmp, &(result.objectStoreId)) {
			return false
		}
	}
	sliceValue = sliceValue[objectStoreIdBytes:]
	{
		tmp := sliceValue[0:indexIdBytes]
		if !DecodeInt(&tmp, &(result.indexId)) {
			return false
		}
	}
	sliceValue = sliceValue[indexIdBytes:]
	*slice = sliceValue
	return true
}
