package leveldbCoding

import (
	"testing"
)

func TestFunctions(t *testing.T) {
	t.Run("CompareGeneric", func(t *testing.T) {
		a := []byte{1, 2, 3, 4}
		b := []byte{1, 2, 3, 4}
		var ok = false
		res := CompareGeneric[DataBaseFreeListKey](a, b, false, &ok)
		if ok == true || res != 0 {
			t.Error("not ok")
		}
	})
}
