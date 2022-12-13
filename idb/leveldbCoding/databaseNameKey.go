package leveldbCoding

type DatabaseNameKey struct {
	origin       U16string
	databaseName U16string
}

func (k DatabaseNameKey) EncodeMinKeyForOrigin(originIdentifier string) string {
	return k.Encode(originIdentifier, U16string{})
}

func (k DatabaseNameKey) EncodeStopKeyForOrigin(originIdentifier string) string {
	return k.EncodeMinKeyForOrigin(originIdentifier + "\x01")
}

func (k DatabaseNameKey) Encode(originIdentifier string, databaseName U16string) string {
	ret := string([]byte{0, 0, 0, 0, KDatabaseNameTypeByte})
	EncodeStringWithLength(ASCIIToUTF16(originIdentifier), &ret)
	EncodeStringWithLength(databaseName, &ret)
	return ret
}

func (k DatabaseNameKey) Decode(slice *[]byte, result *DatabaseNameKey) bool {
	var prefix KeyPrefix
	if !(KeyPrefix{}).Decode(slice, &prefix) {
		return false
	}
	if prefix.databaseId != 0 || prefix.objectStoreId != 0 || prefix.indexId != 0 {
		return false // DCHECK
	}

	var typeByte byte
	if !DecodeByte(slice, &typeByte) {
		return false
	}
	if typeByte != KDatabaseNameTypeByte {
		return false // DCHECK
	}

	if !DecodeStringWithLength(slice, &result.origin) {
		return false
	}
	if !DecodeStringWithLength(slice, &result.databaseName) {
		return false
	}
	return true
}

func (k DatabaseNameKey) Compare(other DatabaseNameKey) int {
	if x := CompareU16String(k.origin, other.origin); x != 0 {
		return x
	}
	return CompareU16String(k.databaseName, other.databaseName)
}
