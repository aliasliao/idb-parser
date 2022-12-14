package databaseNameKey

import (
	"idb-parser/idb/leveldbCoding"
	"idb-parser/idb/leveldbCoding/keyPrefix"
)

type DatabaseNameKey struct {
	Origin       leveldbCoding.U16string
	DatabaseName leveldbCoding.U16string
}

func (k DatabaseNameKey) EncodeMinKeyForOrigin(originIdentifier string) string {
	return k.Encode(originIdentifier, leveldbCoding.U16string{})
}

func (k DatabaseNameKey) EncodeStopKeyForOrigin(originIdentifier string) string {
	return k.EncodeMinKeyForOrigin(originIdentifier + "\x01")
}

func (k DatabaseNameKey) Encode(originIdentifier string, databaseName leveldbCoding.U16string) string {
	ret := string([]byte{0, 0, 0, 0, leveldbCoding.KDatabaseNameTypeByte})
	leveldbCoding.EncodeStringWithLength(leveldbCoding.ASCIIToUTF16(originIdentifier), &ret)
	leveldbCoding.EncodeStringWithLength(databaseName, &ret)
	return ret
}

func (k DatabaseNameKey) Decode(slice *[]byte, result *DatabaseNameKey) bool {
	var prefix keyPrefix.KeyPrefix
	if !(keyPrefix.KeyPrefix{}).Decode(slice, &prefix) {
		return false
	}
	if prefix.DatabaseId != 0 || prefix.ObjectStoreId != 0 || prefix.IndexId != 0 {
		return false // DCHECK
	}

	var typeByte byte
	if !leveldbCoding.DecodeByte(slice, &typeByte) {
		return false
	}
	if typeByte != leveldbCoding.KDatabaseNameTypeByte {
		return false // DCHECK
	}

	if !leveldbCoding.DecodeStringWithLength(slice, &result.Origin) {
		return false
	}
	if !leveldbCoding.DecodeStringWithLength(slice, &result.DatabaseName) {
		return false
	}
	return true
}

func (k DatabaseNameKey) Compare(other DatabaseNameKey) int {
	if x := leveldbCoding.CompareU16String(k.Origin, other.Origin); x != 0 {
		return x
	}
	return leveldbCoding.CompareU16String(k.DatabaseName, other.DatabaseName)
}
