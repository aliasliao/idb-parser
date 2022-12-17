package main

import (
	"log"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"

	"idb-parser/idb/leveldbCoding/compare"
	"idb-parser/idb/metadataCoding"
)

func main() {
	dbPath := "./data/IndexedDB/https_web.haiserve.com_0.indexeddb.leveldb"
	originIdentifier := "https_web.haiserve.com_0@1"

	options := opt.Options{
		Comparer:       compare.Comparator{},
		ErrorIfExist:   false,
		ErrorIfMissing: true,
		ReadOnly:       true,
	}
	db, err := leveldb.OpenFile(dbPath, &options)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	res := metadataCoding.ReadDatabaseNamesAndVersions(db, originIdentifier)
	for _, nv := range res {
		log.Printf("version: %v, name: %s\n", nv.Version, nv.Name.ToString())
	}
}
