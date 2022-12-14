package metadataCoding

import (
	"fmt"
	"log"

	"github.com/syndtr/goleveldb/leveldb"

	"idb-parser/idb/leveldbCoding"
	"idb-parser/idb/leveldbCoding/databaseNameKey"
)

type NameAndVersion struct {
	name    leveldbCoding.U16string
	version int64
}

func ReadDatabaseNamesAndVersions(db *leveldb.DB, originIdentifier string) []NameAndVersion {
	startKey := databaseNameKey.DatabaseNameKey{}.EncodeMinKeyForOrigin(originIdentifier)
	stopKey := databaseNameKey.DatabaseNameKey{}.EncodeStopKeyForOrigin(originIdentifier)
	iter := db.NewIterator(nil, nil)
	found := iter.Seek([]byte(startKey))
	fmt.Printf("it.key=%v\nit.ketStr=%v\nstartKey=%v\nstopKey=%v\nfound=%v\n", iter.Key(), string(iter.Key()), []byte(startKey), stopKey, found)
	iter.Release()
	err := iter.Error()
	if err != nil {
		log.Fatalln(err)
	}
	return []NameAndVersion{}
}
