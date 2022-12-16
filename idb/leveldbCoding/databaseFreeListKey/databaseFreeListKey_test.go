package databaseFreeListKey

import (
	"testing"

	"idb-parser/idb/leveldbCoding/compare"
)

func TestFunctions(t *testing.T) {
	t.Run("CompareGeneric", func(t *testing.T) {
		a := []byte{1, 2, 3, 4}
		b := []byte{1, 2, 3, 4}
		var ok = false
		res := compare.CompareGeneric[DataBaseFreeListKey](a, b, false, &ok)
		if ok == true || res != 0 {
			t.Error("not ok")
		}
	})
}
