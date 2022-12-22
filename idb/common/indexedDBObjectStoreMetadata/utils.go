package indexedDBObjectStoreMetadata

import (
	"idb-parser/idb/common"
	"idb-parser/idb/common/indexedDBIndexMetadata"
	"idb-parser/idb/common/indexedDBKeyPath"
)

type IndexedDBObjectStoreMetadata struct {
	Name          common.U16string
	Id            int64
	KeyPath       indexedDBKeyPath.IndexedDBKeyPath
	AuthIncrement bool
	MaxIndexId    int64
	Indexes       map[int64]indexedDBIndexMetadata.IndexedDBIndexMetadata
}
