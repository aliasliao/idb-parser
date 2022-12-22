package indexedDBKeyPath

import (
	"idb-parser/idb/common"
	"idb-parser/idb/common/mojom/idbKeyPathType"
)

type IndexedDBKeyPath struct {
	Type   idbKeyPathType.IDBKeyPathType
	String common.U16string
	Array  []common.U16string
}
