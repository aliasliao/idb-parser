package leveldbCoding

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
