package databaseNameKey

import "idb-parser/idb/coding"

type DataBaseNameKey struct {
}

func EncodeMinKeyForOrigin(originIdentifier string) string {
	return Encode(originIdentifier, coding.U16string{})
}

func EncodeStopKeyForOrigin(originIdentifier string) string {
	return EncodeMinKeyForOrigin(originIdentifier + "\x01")
}

func Encode(originIdentifier string, databaseName coding.U16string) string {
	ret := string([]byte{0, 0, 0, 0, coding.KDatabaseNameTypeByte})
	coding.EncodeStringWithLength(coding.ASCIIToUTF16(originIdentifier), &ret)
	coding.EncodeStringWithLength(databaseName, &ret)
	return ret
}
