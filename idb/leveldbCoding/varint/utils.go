package varint

func EncodeVarInt(from int64, into *string) {
	n := uint64(from)
	var buf []byte
	if n == 0 {
		buf = []byte{0}
	}
	for n > 0 {
		c := byte(n & 0x7f)
		n >>= 7
		if n > 0 {
			c |= 0x80
		}
		buf = append(buf, c)
	}
	*into += string(buf)
}

func DecodeVarInt(from *[]byte, into *int64) bool {
	fromValue := *from
	var ret uint64 = 0
	i := 0
	shift := 0
	for {
		if i >= len(fromValue) || shift >= 64 {
			return false
		}
		b := fromValue[i]
		ret |= uint64(b&0x7f) << shift
		shift += 7
		i += 1
		if b&0x80 == 0 {
			break
		}
	}
	*into = int64(ret)
	*from = fromValue[i:]
	return true
}
