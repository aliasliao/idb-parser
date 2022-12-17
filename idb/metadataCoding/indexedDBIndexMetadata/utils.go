package indexedDBIndexMetadata

import (
	"idb-parser/idb/leveldbCoding"
	"idb-parser/idb/metadataCoding/indexedDBKeyPath"
)

type IndexedDBIndexMetadata struct {
	Name       leveldbCoding.U16string
	Id         int64
	KeyPath    indexedDBKeyPath.IndexedDBKeyPath
	Unique     bool
	MultiEntry bool
}

const (
	KInvalidId int64 = -1
)
