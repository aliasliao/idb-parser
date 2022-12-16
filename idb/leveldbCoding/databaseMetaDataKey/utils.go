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
