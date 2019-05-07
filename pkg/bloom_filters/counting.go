package filter

import (
	counting "github.com/tylertreat/BoomFilters"
)

// Counting
type Counting struct {
	filter *counting.CountingBloomFilter
}

func NewCounting(n uint) *Counting {
	return &Counting{
		filter: counting.NewCountingBloomFilter(n, 30, 0.001),
	}
}

func (c *Counting) Add(value *[]byte) bool {
	c.filter.Add(*value)
	return true
}

func (c *Counting) Remove(value *[]byte) bool {
	return c.filter.TestAndRemove(*value)
}

func (c *Counting) Test(value *[]byte) bool {
	return c.filter.Test(*value)
}

func Load(buckets []byte, capacity uint, k uint, count uint, indexes []uint) *Counting {
	return &Counting{
		filter: counting.CountingFilterLoadWith(buckets, 30, capacity, k, count, indexes),
	}
}

func (c *Counting) Dump() (b []byte, m uint, k uint, count uint, indexes []uint) {
	b = c.filter.GetBuckets()
	m = c.filter.Capacity()
	k = c.filter.K()
	count = c.filter.Count()
	indexes = c.filter.GetIndexBuffer()
	return
}
