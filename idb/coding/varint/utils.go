package varint

func EncodeVarInt(from int, into *string) {
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
