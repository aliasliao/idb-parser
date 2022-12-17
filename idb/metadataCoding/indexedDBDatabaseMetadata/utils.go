package indexedDBDatabaseMetadata

import (
	"idb-parser/idb/leveldbCoding"
	"idb-parser/idb/metadataCoding/indexedDBObjectStoreMetadata"
)

type IndexedDBDatabaseMetadata struct {
	Name             leveldbCoding.U16string
	Id               int64
	Version          int64
	MaxObjectStoreId int64
	ObjectStores     map[int64]indexedDBObjectStoreMetadata.IndexedDBObjectStoreMetadata
	WasColdOpen      bool
}

const (
	NoVersion      = -1
	DefaultVersion = 0
)
