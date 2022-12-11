package coding

import "idb-parser/idb/keyPrefix"

func Compare(a, b []byte, onlyCompareIndexKeys bool) int {
	sliceA := append([]byte{}, a...)
	sliceB := append([]byte{}, b...)
	prefixA := keyPrefix.KeyPrefix{}
	prefixB := keyPrefix.KeyPrefix{}

	okA := keyPrefix.Decode(&sliceA, &prefixA)
	okB := keyPrefix.Decode(&sliceB, &prefixB)
	if !okA || !okB {
		return 0
	}

	if x := prefixA.Compare(prefixB); x != 0 {
		return x
	}

	switch prefixA.Type() {
	case keyPrefix.GlobalMetadata:

	}

	return 1
}
