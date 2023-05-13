package reply

type PongReply struct{}
type OkReply struct{}
type NullBulkReply struct{}
type EmptyMultiBulkReply struct{}
type NoReply struct{}

var okBytes = []byte("+OK\r\n")
var pongBytes = []byte("+PONG\r\n")
var nullBulkBytes = []byte("$-1\r\n")
var emptyMultiBulkBytes = []byte("*0\r\n")
var noReplyBytes = []byte("")

func (o OkReply) ToBytes() []byte {
	return okBytes
}

func (p PongReply) ToBytes() []byte {
	return pongBytes
}

func (n NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

func (n EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}

func (n NoReply) ToBytes() []byte {
	return noReplyBytes
}

var thePongReply = new(PongReply)
var theNullReply = new(NullBulkReply)
var theOkReply = new(OkReply)
var theEmptyMultiBulkReply = new(EmptyMultiBulkReply)
var theNoReply = new(NoReply)

func MakePongReply() *PongReply {
	return thePongReply
}

func MakeOkReply() *OkReply {
	return theOkReply
}

func MakeNullBulkReply() *NullBulkReply {
	return theNullReply
}

func MakeEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return theEmptyMultiBulkReply
}

func MakeNoReply() *NoReply {
	return theNoReply
}
