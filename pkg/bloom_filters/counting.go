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
		filter: counting.NewCountingBloomFilter(n, 15, 0.01),
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
