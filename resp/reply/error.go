package reply

type UnknownErrReply struct{}
type ArgNumErrReply struct{ Cmd string }
type SyntaxErrReply struct{ Msg string }
type ProtocolErrReply struct{ Msg string }
type WrongTypeErrReply struct{}

var unknownErrorBytes = []byte("-Err unknown\r\n")
var syntaxErrBytes = []byte("-Err syntax error\r\n")
var wrongTypeErrBytes = []byte("-WRONG-TYPE Operation against a key holding the wrong kind of value\r\n")

var theUnknownErrReply = &UnknownErrReply{}

func (u UnknownErrReply) ToBytes() []byte {
	return unknownErrorBytes
}

func (u UnknownErrReply) Error() string {
	return string(wrongTypeErrBytes)
}

func (a ArgNumErrReply) ToBytes() []byte {
	return []byte("-Err wrong number of arguments for '" + a.Cmd + "' command\r\n")
}

func (a ArgNumErrReply) Error() string {
	return "-Err wrong number of arguments for '" + a.Cmd + "' command\r\n"
}

func (s SyntaxErrReply) ToBytes() []byte {
	return syntaxErrBytes
}

func (s SyntaxErrReply) Error() string {
	return "-Err syntax error"
}

func (r *WrongTypeErrReply) ToBytes() []byte {
	return wrongTypeErrBytes
}

func (r *WrongTypeErrReply) Error() string {
	return "WRONG-TYPE Operation against a key holding the wrong kind of value"
}
func (r *ProtocolErrReply) ToBytes() []byte {
	return []byte("-Err Protocol error: '" + r.Msg + "'\r\n")
}

func (r *ProtocolErrReply) Error() string {
	return "-Err Protocol error: '" + r.Msg
}

func MakeSyntaxErrReply(msg string) *SyntaxErrReply {
	return &SyntaxErrReply{msg}
}

func MakeUnknownErrReply() *UnknownErrReply {
	return theUnknownErrReply
}

func MakeProtocolErrReply(msg string) *ProtocolErrReply {
	return &ProtocolErrReply{msg}
}

func MakeArgNumErrReply(cmd string) *ArgNumErrReply {
	return &ArgNumErrReply{
		Cmd: cmd,
	}
}
