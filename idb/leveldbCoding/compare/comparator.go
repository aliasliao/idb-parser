package compare

type Comparator struct {
}

func (c Comparator) Compare(a, b []byte) int {
	return Compare(a, b, false)
}

func (c Comparator) Name() string {
	return "idb_cmp1"
}

func (c Comparator) Separator(dst, a, b []byte) []byte {
	return nil
}

func (c Comparator) Successor(dst, b []byte) []byte {
	return nil
}
