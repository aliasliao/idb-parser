package indexedDBDatabaseMetadata

import (
	"idb-parser/idb/common"
	"idb-parser/idb/common/indexedDBObjectStoreMetadata"
)

type IndexedDBDatabaseMetadata struct {
	Name             common.U16string
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
