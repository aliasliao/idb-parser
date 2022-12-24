package main

import (
	"flag"
	"log"
	"path/filepath"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"

	"idb-parser/idb/common"
	"idb-parser/idb/leveldbCoding/compare"
	"idb-parser/idb/metadataCoding"
)

func getOriginIdentifier(dbPath string) string {
	dirName := filepath.Base(dbPath)
	origin := strings.TrimSuffix(dirName, ".indexeddb.leveldb")
	return origin + "@1"
}

const defaultDbPath = "./data/IndexedDB/https_web.haiserve.com_0.indexeddb.leveldb"

func main() {
	var dbPath = defaultDbPath
	flag.StringVar(&dbPath, "dbPath", defaultDbPath, `leveldb dir path ending with ".indexeddb.leveldb"`)
	flag.Parse()

	originIdentifier := getOriginIdentifier(dbPath)

	options := opt.Options{
		Comparer:       compare.Comparator{},
		ErrorIfExist:   false,
		ErrorIfMissing: true,
		ReadOnly:       true,
	}
	db, err := leveldb.OpenFile(dbPath, &options)
	if err != nil {
		log.Fatalln("fail to open file:", err)
	}
	defer func(db *leveldb.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(db)

	namesAndVersions, err := metadataCoding.ReadDatabaseNamesAndVersions(db, originIdentifier)
	if err != nil || len(*namesAndVersions) == 0 {
		log.Fatalf("ReadDatabaseNamesAndVersions error: %v\n", err)
	}
	for _, nv := range *namesAndVersions {
		log.Printf("id: %v, name: %v, version: %2v\n", nv.Id, nv.Name.ToString(), nv.Version)

		if metadata, err := metadataCoding.ReadMetadataForDatabaseName(db, originIdentifier, nv.Name); err != nil {
			log.Fatalf("ReadMetadataForDatabaseName error: %v\n", err)
		} else {
			log.Println()
			log.Printf("objectStores length: %v, maxObjectStoreId: %v\n", len(metadata.ObjectStores), metadata.MaxObjectStoreId)
			for objectStoreId, objectStoreMetadata := range metadata.ObjectStores {
				log.Println()
				log.Printf("  objectStoreId: %v\n", objectStoreId)
				log.Printf("  name: %s\n", objectStoreMetadata.Name.ToString())
				log.Printf("  keyPath: type=%v, string=%s, array=%v\n", objectStoreMetadata.KeyPath.Type, objectStoreMetadata.KeyPath.String.ToString(), common.StringifyArray(objectStoreMetadata.KeyPath.Array))
				log.Println()
				log.Printf("  indexes length: %v, maxIndexId: %v\n", len(objectStoreMetadata.Indexes), objectStoreMetadata.MaxIndexId)
				for indexId, indexMetadata := range objectStoreMetadata.Indexes {
					log.Printf("    indexId: %v\n", indexId)
					log.Printf("    name: %v\n", indexId)
					log.Printf("    keyPath: type=%v, string=%s, array=%v\n", indexMetadata.KeyPath.Type, indexMetadata.KeyPath.String.ToString(), common.StringifyArray(indexMetadata.KeyPath.Array))
					log.Printf("    unique: %v\n", indexMetadata.Unique)
					log.Printf("    multiEntry: %v\n", indexMetadata.MultiEntry)
					log.Println()
				}
				log.Println()
			}
		}
		log.Println()
	}
}
