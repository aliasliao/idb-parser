package coding

import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

var nativeEndian binary.ByteOrder

func init() {
	buf := [2]byte{}
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)

	switch buf {
	case [2]byte{0xCD, 0xAB}:
		nativeEndian = binary.LittleEndian
		fmt.Println("is little endian")
	case [2]byte{0xAB, 0xCD}:
		nativeEndian = binary.BigEndian
		fmt.Println("is big endian")
	default:
		panic("Could not determine native endianness.")
	}
}
