package main

import (
    "fmt"
    "log"

    "github.com/syndtr/goleveldb/leveldb"
    "github.com/syndtr/goleveldb/leveldb/opt"
)

func main() {
    dbPath := "./data/IndexedDB/https_web.haiserve.com_0.indexeddb.leveldb"

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
        fmt.Printf("key=%v\nvalue=%v\n\n", string(key), string(value))
    }
    iter.Release()
    err = iter.Error()
    if err != nil {
        log.Fatalln(err)
    }
}
