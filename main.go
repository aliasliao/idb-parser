package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"

	"idb-parser/idb/metadataCoding"
)

func main() {
	dbPath := "./data/IndexedDB/https_web.haiserve.com_0.indexeddb.leveldb"
	originIdentifier := "https_web.haiserve.com_0@1"

	options := opt.Options{
		Comparer:       bytesComparer{},
		ErrorIfExist:   false,
		ErrorIfMissing: true,
		ReadOnly:       true,
	}
	db, err := leveldb.OpenFile(dbPath, &options)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		key := iter.Key()
		value := iter.Value()
		if bytes.Contains(key, []byte{0, 0, 0, 0, 201}) {
			fmt.Printf("key=%v\nvalue=%v\n,keyBytes=%v\n\n", string(key), string(value), key)
		}
	}
	iter.Release()
	err = iter.Error()
	if err != nil {
		log.Fatalln(err)
	}

	res := metadataCoding.ReadDatabaseNamesAndVersions(db, originIdentifier)
	for _, nv := range res {
		log.Printf("version: %v, name: %s\n", nv.Version, nv.Name)
	}
}
