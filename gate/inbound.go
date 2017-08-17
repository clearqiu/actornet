package gate

import (
	"github.com/davyxu/actornet/actor"
	"github.com/davyxu/actornet/proto"
	"github.com/davyxu/cellnet"
)

type inboundHandler struct {
}

// 收到客户端发过来的消息
func (self *inboundHandler) Call(ev *cellnet.Event) {

	if ev.Type == cellnet.Event_Recv {

		tag := ev.Ses.Tag()
		if tag != nil {

			backendPID := tag.(*actor.PID)

			outboundPID := PIDBySession(ev.Ses)

			log.Debugf("direct route: %s -> %s (%s)", outboundPID.String(), backendPID.String())

			backendPID.NotifyBySender(ev.Msg, outboundPID)

		} else {

			// TODO 检查消息是否在透传列表中

			switch ev.Msg.(type) {
			case *proto.BindClientREQ:
				backendAssit.NotifyBySender(&proto.BindClientREQ{ev.Ses.ID()}, receiptor)
			}

		}

	}
}

func newInboundHandler() cellnet.EventHandler {
	return &inboundHandler{}
}

// 后台的辅助actor
var backendAssit *actor.PID

func init() {

	backendAssit = actor.NewPID("server", "gate_assit")
}
