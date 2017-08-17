package nexus

import (
	"github.com/davyxu/actornet/proto"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/socket"
)

// 启动本机的listen
func Connect(address string, domain string) {

	peer := socket.NewConnector(nil)
	peer.Start(address)

	cellnet.RegisterMessage(peer, "coredef.SessionConnected", func(ev *cellnet.Event) {

		ev.Send(&proto.DomainIdentifyACK{
			Domain: domain,
		})
	})

	cellnet.RegisterMessage(peer, "proto.DomainIdentifyACK", func(ev *cellnet.Event) {
		msg := ev.Msg.(*proto.DomainIdentifyACK)

		addServiceSession(msg.Domain, ev.Ses)
	})

	cellnet.RegisterMessage(peer, "coredef.SessionClosed", func(ev *cellnet.Event) {

		removeServiceSession(ev.Ses)

	})

	cellnet.RegisterMessage(peer, "proto.RouteACK", onRouter)
}
