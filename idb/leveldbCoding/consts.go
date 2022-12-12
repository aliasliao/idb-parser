package leveldbCoding

const (
	KDefaultInlineBufferSize uint32 = 32

	KIndexedDBKeyNullTypeByte   byte = 0
	KIndexedDBKeyStringTypeByte byte = 1
	KIndexedDBKeyDateTypeByte   byte = 2
	KIndexedDBKeyNumberTypeByte byte = 3
	KIndexedDBKeyArrayTypeByte  byte = 4
	KIndexedDBKeyMinKeyTypeByte byte = 5
	KIndexedDBKeyBinaryTypeByte byte = 6

	KIndexedDBKeyPathTypeCodedByte1 byte = 0
	KIndexedDBKeyPathTypeCodedByte2 byte = 0

	KIndexedDBKeyPathNullTypeByte   byte = 0
	KIndexedDBKeyPathStringTypeByte byte = 1
	KIndexedDBKeyPathArrayTypeByte  byte = 2

	KObjectStoreDataIndexId byte = 1
	KExistsEntryIndexId     byte = 2
	KBlobEntryIndexId       byte = 3

	KSchemaVersionTypeByte           byte = 0
	KMaxDatabaseIdTypeByte           byte = 1
	KDataVersionTypeByte             byte = 2
	KRecoveryBlobJournalTypeByte     byte = 3
	KActiveBlobJournalTypeByte       byte = 4
	KEarliestSweepTimeTypeByte       byte = 5
	KEarliestCompactionTimeTypeByte  byte = 6
	KMaxSimpleGlobalMetaDataTypeByte byte = 7 // Insert before this and increment.
	KScopesPrefixByte                byte = 50
	KDatabaseFreeListTypeByte        byte = 100
	KDatabaseNameTypeByte            byte = 201

	KObjectStoreMetaDataTypeByte byte = 50
	KIndexMetaDataTypeByte       byte = 100
	KObjectStoreFreeListTypeByte byte = 150
	KIndexFreeListTypeByte       byte = 151
	KObjectStoreNamesTypeByte    byte = 200
	KIndexNamesKeyTypeByte       byte = 201

	KObjectMetaDataTypeMaximum byte = 255
	KIndexMetaDataTypeMaximum  byte = 255

	KDatabaseLockPartition    int = 0
	KObjectStoreLockPartition int = 1
)
