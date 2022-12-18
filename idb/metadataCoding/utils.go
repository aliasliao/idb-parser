package metadataCoding

import (
	"fmt"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"

	"idb-parser/idb/leveldbCoding"
	"idb-parser/idb/leveldbCoding/compare"
	"idb-parser/idb/leveldbCoding/databaseMetaDataKey"
	"idb-parser/idb/leveldbCoding/databaseNameKey"
	"idb-parser/idb/leveldbCoding/keyPrefix"
	"idb-parser/idb/leveldbCoding/mojom/idbKeyPathType"
	"idb-parser/idb/leveldbCoding/objectStoreMetaDataKey"
	"idb-parser/idb/leveldbCoding/varint"
	"idb-parser/idb/metadataCoding/indexedDBDatabaseMetadata"
	"idb-parser/idb/metadataCoding/indexedDBKeyPath"
	"idb-parser/idb/metadataCoding/indexedDBObjectStoreMetadata"
)

type NameAndVersion struct {
	Name    leveldbCoding.U16string
	Id      int64
	Version int64
}

func GetVarInt(db *leveldb.DB, key *string) (*int64, error) {
	slice := []byte(*key)
	value, err := db.Get(slice, nil)
	if err != nil {
		return nil, err
	}
	var foundInt int64 = 0
	if !varint.DecodeVarInt(&value, &foundInt) || len(value) != 0 {
		return nil, fmt.Errorf("fail to DecodeVarInt")
	}
	return &foundInt, nil
}

func GetInt(db *leveldb.DB, key *string) (*int64, error) {
	slice := []byte(*key)
	value, err := db.Get(slice, nil)
	if err != nil {
		return nil, err
	}
	var foundInt int64 = 0
	if !leveldbCoding.DecodeInt(&value, &foundInt) || len(value) != 0 {
		return nil, fmt.Errorf("fail to DecodeInt")
	}
	return &foundInt, nil
}

func GetMaxObjectStoreId(db *leveldb.DB, databaseId int64) (*int64, error) {
	key := databaseMetaDataKey.DatabaseMetaDataKey{}.Encode(databaseId, databaseMetaDataKey.MaxObjectStoreId)
	if maxObjectStoreId, err := GetInt(db, &key); err != nil {
		return nil, fmt.Errorf("fail to get max object store id: %w", err)
	} else {
		if *maxObjectStoreId < 0 {
			return nil, fmt.Errorf("maxObjectStoreId < 0")
		}
		return maxObjectStoreId, nil
	}
}

func CheckObjectStoreAndMetaDataType(it iterator.Iterator, stopKey string, objectStoreId int64, metaDataType objectStoreMetaDataKey.MetaDataType) bool {
	if !it.Valid() || compare.CompareKeys(it.Key(), []byte(stopKey)) >= 0 {
		return false
	}
	slice := it.Key()
	var metaDataKey objectStoreMetaDataKey.ObjectStoreMetaDataKey
	if !(objectStoreMetaDataKey.ObjectStoreMetaDataKey{}).Decode(&slice, &metaDataKey) || len(slice) != 0 {
		panic("fail to decode objectStoreMetaDataKey")
	}
	if metaDataKey.ObjectStoreId != objectStoreId || metaDataKey.MetaDataType != metaDataType {
		return false
	}
	return true
}

func ReadObjectStores(db *leveldb.DB, databaseId int64) (*map[int64]indexedDBObjectStoreMetadata.IndexedDBObjectStoreMetadata, error) {
	if !keyPrefix.IsValidDatabaseId(databaseId) {
		return nil, fmt.Errorf("invalid databaseId")
	}

	startKey := objectStoreMetaDataKey.ObjectStoreMetaDataKey{}.Encode(databaseId, 1, objectStoreMetaDataKey.Name)
	stopKey := objectStoreMetaDataKey.ObjectStoreMetaDataKey{}.EncodeMaxKey(databaseId)

	var objectStores map[int64]indexedDBObjectStoreMetadata.IndexedDBObjectStoreMetadata

	it := db.NewIterator(nil, nil)
	ok := it.Seek([]byte(startKey))
	for ok && it.Valid() && compare.CompareKeys(it.Key(), []byte(stopKey)) < 0 {
		var metaDataKey objectStoreMetaDataKey.ObjectStoreMetaDataKey
		{
			slice := it.Key()
			if !(objectStoreMetaDataKey.ObjectStoreMetaDataKey{}).Decode(&slice, &metaDataKey) || len(slice) != 0 {
				panic("fail to decode metaDataKey")
			}
			if metaDataKey.MetaDataType != objectStoreMetaDataKey.Name {
				ok = it.Next()
				continue
			}
		}

		objectStoreId := metaDataKey.ObjectStoreId
		var objectStoreName leveldbCoding.U16string
		{
			slice := it.Value()
			if !leveldbCoding.DecodeString(&slice, &objectStoreName) || len(slice) != 0 {
				panic("fail to decode objectStoreName")
			}
		}

		ok = it.Next()
		if !ok || !CheckObjectStoreAndMetaDataType(it, stopKey, objectStoreId, objectStoreMetaDataKey.KeyPath) {
			break
		}
		var keyPath indexedDBKeyPath.IndexedDBKeyPath
		{
			slice := it.Value()
			if !DecodeIDBKeyPath(&slice, &keyPath) || len(slice) != 0 {
				panic("fail to decode IDBKeyPath")
			}
		}

		ok = it.Next()
		if !ok || !CheckObjectStoreAndMetaDataType(it, stopKey, objectStoreId, objectStoreMetaDataKey.AutoIncrement) {
			break
		}
		var autoIncrement bool
		{
			slice := it.Value()
			if !leveldbCoding.DecodeBool(&slice, &autoIncrement) || len(slice) != 0 {
				panic("fail to decode autoIncrement")
			}
		}

		ok = it.Next()
		if !ok || !CheckObjectStoreAndMetaDataType(it, stopKey, objectStoreId, objectStoreMetaDataKey.Evictable) {
			break
		}

		ok = it.Next()
		if !ok || !CheckObjectStoreAndMetaDataType(it, stopKey, objectStoreId, objectStoreMetaDataKey.LastVersion) {
			break
		}

		ok = it.Next()
		if !ok || !CheckObjectStoreAndMetaDataType(it, stopKey, objectStoreId, objectStoreMetaDataKey.MaxIndexId) {
			break
		}
		var maxIndexId int64
		{
			slice := it.Value()
			if !leveldbCoding.DecodeInt(&slice, &maxIndexId) || len(slice) != 0 {
				panic("fail to decode maxIndexId")
			}
		}

		ok = it.Next()
		if !ok {
			break
		}
		if CheckObjectStoreAndMetaDataType(it, stopKey, objectStoreId, objectStoreMetaDataKey.HasKeyPath) {
			var hasKeyPath bool
			{
				slice := it.Value()
				if !leveldbCoding.DecodeBool(&slice, &hasKeyPath) {
					panic("fail to decode hasKeyPath")
				}
			}
			if !hasKeyPath && keyPath.Type == idbKeyPathType.String && len(keyPath.String) != 0 {
				break
			}
			if !hasKeyPath {
				keyPath = indexedDBKeyPath.IndexedDBKeyPath{Type: idbKeyPathType.Null}
			}
			if ok = it.Next(); !ok {
				break
			}
		}

		var keyGeneratorCurrentNumber int64 = -1
		if CheckObjectStoreAndMetaDataType(it, stopKey, objectStoreId, objectStoreMetaDataKey.KeyGeneratorCurrentNumber) {
			slice := it.Value()
			if !leveldbCoding.DecodeInt(&slice, &keyGeneratorCurrentNumber) || len(slice) != 0 {
				panic("fail to decode keyGeneratorCurrentNumber")
			}
			if keyGeneratorCurrentNumber < objectStoreMetaDataKey.KKeyGeneratorInitialNumber {
				panic("keyGeneratorCurrentNumber < objectStoreMetaDataKey.KKeyGeneratorInitialNumber")
			}
			if ok = it.Next(); !ok {
				break
			}
		}

		metadata := indexedDBObjectStoreMetadata.IndexedDBObjectStoreMetadata{
			Name:          objectStoreName,
			Id:            objectStoreId,
			KeyPath:       keyPath,
			AuthIncrement: autoIncrement,
			MaxIndexId:    maxIndexId,
			Indexes:       nil,
		}
		if !ReadIndexes(db, databaseId, objectStoreId, &metadata.Indexes) {
			break
		}
		objectStores[objectStoreId] = metadata
	}

	it.Release()
	err := it.Error()
	if err != nil {
		return nil, err
	}

	return &objectStores, nil
}

func ReadDatabaseNamesAndVersions(db *leveldb.DB, originIdentifier string) (*[]NameAndVersion, error) {
	var ret []NameAndVersion

	startKey := databaseNameKey.DatabaseNameKey{}.EncodeMinKeyForOrigin(originIdentifier)
	stopKey := databaseNameKey.DatabaseNameKey{}.EncodeStopKeyForOrigin(originIdentifier)

	it := db.NewIterator(nil, nil)
	ok := it.Seek([]byte(startKey))
	for ok && it.Valid() && compare.CompareKeys(it.Key(), []byte(stopKey)) < 0 {
		// Decode database Name (in iterator key).
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

		// Look up Version by id.
		dbVersion := int64(indexedDBDatabaseMetadata.DefaultVersion)
		metaDataKey := databaseMetaDataKey.DatabaseMetaDataKey{}.Encode(dbId, databaseMetaDataKey.UserVersion)
		if foundInt, err := GetVarInt(db, &metaDataKey); err != nil {
			log.Printf("fail to get databaseVersion: %v\n", err)
			continue
		} else {
			dbVersion = *foundInt
		}

		if dbVersion != indexedDBDatabaseMetadata.DefaultVersion {
			ret = append(ret, NameAndVersion{
				Name:    dbNameKey.DatabaseName,
				Id:      dbId,
				Version: dbVersion,
			})
		}

		ok = it.Next()
	}

	it.Release()
	err := it.Error()
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func ReadMetadataForDatabaseName(db *leveldb.DB, originIdentifier string, name leveldbCoding.U16string) (*indexedDBDatabaseMetadata.IndexedDBDatabaseMetadata, error) {
	metadata := indexedDBDatabaseMetadata.IndexedDBDatabaseMetadata{
		Name:    name,
		Version: indexedDBDatabaseMetadata.DefaultVersion,
	}
	dbNameKey := databaseNameKey.DatabaseNameKey{}.Encode(originIdentifier, name)
	if id, err := GetInt(db, &dbNameKey); err != nil {
		return nil, fmt.Errorf("fail to get db id: %w", err)
	} else {
		metadata.Id = *id
	}

	versionKey := databaseMetaDataKey.DatabaseMetaDataKey{}.Encode(metadata.Id, databaseMetaDataKey.UserVersion)
	if version, err := GetVarInt(db, &versionKey); err != nil {
		return nil, fmt.Errorf("fail to get db version: %w", err)
	} else {
		metadata.Version = *version
	}

	if maxObjectStoreId, err := GetMaxObjectStoreId(db, metadata.Id); err != nil {
		return nil, fmt.Errorf("fail to get maxObjectStoreId: %w", err)
	} else {
		metadata.MaxObjectStoreId = *maxObjectStoreId
	}

	blobNumberKey := databaseMetaDataKey.DatabaseMetaDataKey{}.Encode(metadata.Id, databaseMetaDataKey.BlobKeyGeneratorCurrentNumber)
	if currentBlobNumber, err := GetVarInt(db, &blobNumberKey); err != nil {
		return nil, fmt.Errorf("fail to get blob current number: %w", err)
	} else if !databaseMetaDataKey.IsValidBlobNumber(*currentBlobNumber) {
		return nil, fmt.Errorf("blob number not valid")
	}

	if objectStores, err := ReadObjectStores(db, metadata.Id); err != nil {
		return nil, fmt.Errorf("fail to read object stores: %w", err)
	} else {
		metadata.ObjectStores = *objectStores
	}

	return &metadata, nil
}
