package dict

type Consumer func(key string, val any) bool
type Dict interface {
	Get(key string) (val interface{}, exists bool)
	Len() int
	Set(key string, val interface{}) (result int)
	SetIfAbsent(key string, val interface{}) (result int)
	SetIfExists(key string, val interface{}) (result int)
	Remove(key string) (result int)
	ForEach(consumer Consumer)
	Keys() []string
	RandomKeys(count int) []string // values are not the same as each other
	Clear() (result int)
}
