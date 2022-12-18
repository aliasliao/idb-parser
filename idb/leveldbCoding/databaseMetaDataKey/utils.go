package databaseMetaDataKey

import "idb-parser/idb/leveldbCoding/keyPrefix"

type DatabaseMetaDataKey struct {
}

type MetaDataType uint8

const (
	OriginName MetaDataType = iota
	DatabaseName
	UserStringVersion // Obsolete
	MaxObjectStoreId
	UserVersion
	BlobKeyGeneratorCurrentNumber
	MaxSimpleMetadataType
)

const (
	KAllBlobsNumber                   int64 = 1
	KBlobNumberGeneratorInitialNumber int64 = 2
	KInvalidBlobNumber                int64 = -1
)

func (k DatabaseMetaDataKey) Encode(databaseId int64, metaDataType MetaDataType) string {
	prefix := keyPrefix.KeyPrefix{
		DatabaseId:    databaseId,
		ObjectStoreId: 0,
		IndexId:       0,
	}
	ret := prefix.Encode()
	ret += string(metaDataType)
	return ret
}

func IsValidBlobNumber(blobNumber int64) bool {
	return blobNumber >= KBlobNumberGeneratorInitialNumber
}
