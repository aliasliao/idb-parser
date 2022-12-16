package mojom

type IDBKeyType int32

const (
	Invalid IDBKeyType = iota
	Array
	Binary
	String
	Date
	Number
	None
	Min
	KMinValue IDBKeyType = 0
	KMaxValue IDBKeyType = 7
)
