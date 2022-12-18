package idbKeyPathType

type IDBKeyPathType int32

const (
	Null IDBKeyPathType = iota
	String
	Array

	kMinValue IDBKeyPathType = 0
	kMaxValue IDBKeyPathType = 1
)
