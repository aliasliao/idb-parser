package databaseNameKey

import "idb-parser/idb/leveldbCoding"

type DataBaseNameKey struct {
}

func EncodeMinKeyForOrigin(originIdentifier string) string {
	return Encode(originIdentifier, leveldbCoding.U16string{})
}

func EncodeStopKeyForOrigin(originIdentifier string) string {
	return EncodeMinKeyForOrigin(originIdentifier + "\x01")
}

func Encode(originIdentifier string, databaseName leveldbCoding.U16string) string {
	ret := string([]byte{0, 0, 0, 0, leveldbCoding.KDatabaseNameTypeByte})
	leveldbCoding.EncodeStringWithLength(leveldbCoding.ASCIIToUTF16(originIdentifier), &ret)
	leveldbCoding.EncodeStringWithLength(databaseName, &ret)
	return ret
}
