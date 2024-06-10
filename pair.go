package quickhttp

type pair struct {
	key   []byte
	value buf[int32]
}

func appendPair(p *[]pair, key []byte) {
	if cap(*p) > len(*p) {
		*p = (*p)[:len(*p)+1]
	} else {
		*p = append(*p, pair{})
	}
	newKey := (*p)[len(*p)-1].key
	newKey = append(newKey, key...)
	(*p)[len(*p)-1].key = newKey

	return
}

func setLastValue(p *[]pair, start, end int32, val []byte) {
	last := &(*p)[len(*p)-1]
	last.value.setPosOrBytes(start, end, val, minBufLimit)
}
