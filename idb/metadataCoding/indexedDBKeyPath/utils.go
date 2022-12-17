package indexedDBKeyPath

import "idb-parser/idb/leveldbCoding"

type IndexedDBKeyPath struct {
	Type   IDBKeyPathType
	String leveldbCoding.U16string
	Array  []leveldbCoding.U16string
}

type IDBKeyPathType int32

const (
	Null IDBKeyPathType = iota
	String
	Array

	kMinValue IDBKeyPathType = 0
	kMaxValue IDBKeyPathType = 1
)
