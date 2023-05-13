package resp

// Connection denote a redis connection
type Connection interface {
	Write([]byte) error
	GetDBIndex() int
	SelectDB(int)
}
