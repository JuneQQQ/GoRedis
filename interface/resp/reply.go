package resp

// Reply the data structure sent by the server or client
type Reply interface {
	ToBytes() []byte
}
