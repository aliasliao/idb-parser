package indexedDBObjectStoreMetadata

import (
	"idb-parser/idb/leveldbCoding"
	"idb-parser/idb/metadataCoding/indexedDBIndexMetadata"
	"idb-parser/idb/metadataCoding/indexedDBKeyPath"
)

type IndexedDBObjectStoreMetadata struct {
	Name          leveldbCoding.U16string
	Id            int64
	KeyPath       indexedDBKeyPath.IndexedDBKeyPath
	AuthIncrement bool
	MaxIndexId    int64
	Indexes       map[int64]indexedDBIndexMetadata.IndexedDBIndexMetadata
}
