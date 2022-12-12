package varint

import (
	"bytes"
	"testing"
)

func TestFunctions(t *testing.T) {
	t.Run("EncodeVarInt", func(t *testing.T) {
		t.Run("lt 7 bits", func(t *testing.T) {
			into := string([]byte{0, 0, 0, 0, 201})
			EncodeVarInt(0x12, &into)
			if into != string([]byte{0, 0, 0, 0, 201, 0x12}) {
				t.Error("failed")
			}
		})
		t.Run("eq 7 bits", func(t *testing.T) {
			into := string([]byte{0, 0, 0, 0, 201})
			EncodeVarInt(0x7f, &into)
			if into != string([]byte{0, 0, 0, 0, 201, 0x7f}) {
				t.Errorf("failed, into=%v\n", []byte(into))
			}
		})
		t.Run("gt 7 bits", func(t *testing.T) {
			into := string([]byte{0, 0, 0, 0, 201})
			EncodeVarInt(0x127f, &into)
			if into != string([]byte{0, 0, 0, 0, 201, 0xff, 0x24}) {
				t.Errorf("failed, into=%v\n", []byte(into))
			}
		})
	})

	t.Run("DecodeVarInt", func(t *testing.T) {
		t.Run("lt 7  bits", func(t *testing.T) {
			into := int64(0)
			from := []byte{0x12}
			ok := DecodeVarInt(&from, &into)
			if !ok {
				t.Error("not ok")
			}
			if into != 0x12 || bytes.Compare(from, []byte{}) != 0 {
				t.Errorf("failed, into=0x%x\n", into)
			}
		})
		t.Run("eq 7  bits", func(t *testing.T) {
			into := int64(0)
			from := []byte{0x7f}
			ok := DecodeVarInt(&from, &into)
			if !ok {
				t.Error("not ok")
			}
			if into != 0x7f || bytes.Compare(from, []byte{}) != 0 {
				t.Errorf("failed, into=0x%x\n", into)
			}
		})
		t.Run("gt 7  bits", func(t *testing.T) {
			into := int64(0)
			from := []byte{0xff, 0x24}
			ok := DecodeVarInt(&from, &into)
			if !ok {
				t.Error("not ok")
			}
			if into != 0x127f || bytes.Compare(from, []byte{}) != 0 {
				t.Errorf("failed, into=0x%x\n", into)
			}
		})
	})
}
