package filter

type Filter interface {
	Add(value *[]byte) bool
	Remove(value *[]byte) bool
	Test(value *[]byte) bool
}
