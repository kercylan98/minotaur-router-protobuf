package protobufrouter

import "google.golang.org/protobuf/proto"

type Reader func(message proto.Message)

func (r Reader) ReadTo(message proto.Message) {
	r(message)
}
