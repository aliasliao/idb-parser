package idb

import (
	"fmt"
	"log"

	"github.com/syndtr/goleveldb/leveldb"

	"idb-parser/idb/coding"
	"idb-parser/idb/databaseNameKey"
)

type NameAndVersion struct {
	name    coding.U16string
	version int64
}

func ReadDatabaseNamesAndVersions(db *leveldb.DB, originIdentifier string) []NameAndVersion {
	startKey := databaseNameKey.EncodeMinKeyForOrigin(originIdentifier)
	stopKey := databaseNameKey.EncodeStopKeyForOrigin(originIdentifier)
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
