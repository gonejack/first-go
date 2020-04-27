package util

import "sync/atomic"

type indexer struct {
	count uint64
}

func (idx *indexer) getIndex(len int) int {
	return int(atomic.AddUint64(&idx.count, 1) % uint64(len))
}
