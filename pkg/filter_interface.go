package filter

import "hash"

type Filter interface {
	Add(value *[]byte) bool
	Remove(value *[]byte) bool
	Test(value *[]byte) bool
	Dump() (b []byte, h hash.Hash64, m uint, k uint, count uint, indexes []uint)
}
