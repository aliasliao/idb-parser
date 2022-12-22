package indexedDBIndexMetadata

import (
	"idb-parser/idb/common"
	"idb-parser/idb/common/indexedDBKeyPath"
)

type IndexedDBIndexMetadata struct {
	Name       common.U16string
	Id         int64
	KeyPath    indexedDBKeyPath.IndexedDBKeyPath
	Unique     bool
	MultiEntry bool
}

const (
	KInvalidId int64 = -1
)
