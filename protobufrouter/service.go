package protobufrouter

import (
	"github.com/kercylan98/minotaur/server"
	"github.com/kercylan98/minotaur/server/router"
	"google.golang.org/protobuf/proto"
)

func NewDefault() *Service[MessageID, *Message, *server.Conn] {
	return New[MessageID, *Message, *server.Conn](
		func() *Message {
			return &Message{}
		},
		func(service *Service[MessageID, *Message, *server.Conn], conn *server.Conn) *server.Conn {
			return conn
		},
		func(packet []byte) Reader {
			return func(message proto.Message) {
				if err := proto.Unmarshal(packet, message); err != nil {
					panic(err)
				}
			}
		},
	)
}

func New[MessageID comparable, Message proto.Message, Entity any](
	generateMessageHandler func() Message,
	getEntityHandler func(service *Service[MessageID, Message, Entity], conn *server.Conn) Entity,
	unmarshalHandler func(packet []byte) Reader,
) *Service[MessageID, Message, Entity] {
	s := &Service[MessageID, Message, Entity]{
		Multistage:       router.NewMultistage[HandlerFunc[MessageID, Message, Entity]](),
		generateMessage:  generateMessageHandler,
		getEntity:        getEntityHandler,
		unmarshalHandler: unmarshalHandler,
	}
	return s
}

type HandlerFunc[MessageID comparable, Message proto.Message, Entity any] func(service *Service[MessageID, Message, Entity], entity Entity, reader Reader)

type Service[MessageID comparable, Message proto.Message, Entity any] struct {
	*router.Multistage[HandlerFunc[MessageID, Message, Entity]]
	srv              *server.Server
	generateMessage  func() Message
	getEntity        func(service *Service[MessageID, Message, Entity], conn *server.Conn) Entity
	unmarshalHandler func(packet []byte) Reader
}

func (s *Service[MessageID, Message, Entity]) OnInit(srv *server.Server) {
	s.srv = srv
	srv.RegConnectionReceivePacketEvent(s.OnConnectionReceivePacket)
}

func (s *Service[MessageID, Message, Entity]) Server() *server.Server {
	return s.srv
}

func (s *Service[MessageID, Message, Entity]) OnConnectionReceivePacket(srv *server.Server, conn *server.Conn, packet []byte) {
	var message = s.generateMessage()
	if err := proto.Unmarshal(packet, message); err != nil {
		return
	}

	handler := s.Match(message)
	if handler == nil {
		return
	}

	handler(s, s.getEntity(s, conn), s.unmarshalHandler(packet))
}
