package leveldbCoding

type DataBaseNameKey struct {
}

func EncodeMinKeyForOrigin(originIdentifier string) string {
	return Encode(originIdentifier, U16string{})
}

func EncodeStopKeyForOrigin(originIdentifier string) string {
	return EncodeMinKeyForOrigin(originIdentifier + "\x01")
}

func Encode(originIdentifier string, databaseName U16string) string {
	ret := string([]byte{0, 0, 0, 0, KDatabaseNameTypeByte})
	EncodeStringWithLength(ASCIIToUTF16(originIdentifier), &ret)
	EncodeStringWithLength(databaseName, &ret)
	return ret
}
