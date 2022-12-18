package indexedDBKeyPath

import (
	"idb-parser/idb/leveldbCoding"
	"idb-parser/idb/leveldbCoding/mojom/idbKeyPathType"
)

type IndexedDBKeyPath struct {
	Type   idbKeyPathType.IDBKeyPathType
	String leveldbCoding.U16string
	Array  []leveldbCoding.U16string
}
