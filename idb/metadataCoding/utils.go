package metadataCoding

import (
	"log"

	"github.com/syndtr/goleveldb/leveldb"

	"idb-parser/idb/leveldbCoding"
	"idb-parser/idb/leveldbCoding/databaseMetaDataKey"
	"idb-parser/idb/leveldbCoding/databaseNameKey"
	"idb-parser/idb/leveldbCoding/varint"
)

type NameAndVersion struct {
	name    leveldbCoding.U16string
	version int64
}

func GetVarInt(db *leveldb.DB, key []byte, foundInt *int64) bool {
	result, err := db.Get(key, nil)
	if err != nil {
		return false
	}
	if !varint.DecodeVarInt(&result, foundInt) || len(result) != 0 {
		return false
	}
	return true
}

func ReadDatabaseNamesAndVersions(db *leveldb.DB, originIdentifier string) []NameAndVersion {
	var ret []NameAndVersion

	startKey := databaseNameKey.DatabaseNameKey{}.EncodeMinKeyForOrigin(originIdentifier)
	stopKey := databaseNameKey.DatabaseNameKey{}.EncodeStopKeyForOrigin(originIdentifier)

	it := db.NewIterator(nil, nil)
	ok := it.Seek([]byte(startKey))
	for ok && it.Valid() && leveldbCoding.CompareKeys(it.Key(), []byte(stopKey)) < 0 {
		// Decode database name (in iterator key).
		slice := it.Key()
		var dbNameKey databaseNameKey.DatabaseNameKey
		if !(databaseNameKey.DatabaseNameKey{}).Decode(&slice, &dbNameKey) || len(slice) != 0 {
			log.Println("error getting databaseNameKey")
			continue
		}

		// Decode database id (in iterator value).
		var dbId int64
		valueSlice := it.Value()
		if !leveldbCoding.DecodeInt(&valueSlice, &dbId) || len(valueSlice) != 0 {
			log.Println("error getting databaseId")
			continue
		}

		// Look up version by id.
		dbVersion := int64(DefaultVersion)
		metaDataKey := databaseMetaDataKey.DatabaseMetaDataKey{}.Encode(dbId, databaseMetaDataKey.UserVersion)
		metaDataKeySlice := []byte(metaDataKey)
		if !GetVarInt(db, metaDataKeySlice, &dbVersion) {
			log.Println("error getting databaseVersion")
			continue
		}

		if dbVersion != DefaultVersion {
			ret = append(ret, NameAndVersion{
				name:    dbNameKey.DatabaseName,
				version: dbVersion,
			})
		}

		ok = it.Next()
	}

	it.Release()
	err := it.Error()
	if err != nil {
		log.Fatalln(err)
	}

	return ret
}
