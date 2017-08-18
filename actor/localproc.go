package actor

import (
	"github.com/davyxu/actornet/mailbox"
	"github.com/davyxu/actornet/proto"
	"github.com/davyxu/actornet/util"
)

type Process interface {
	Tell(interface{})

	Stop()

	PID() *PID
}

type localProcess struct {
	mailbox mailbox.MailBox

	pid *PID

	a Actor
}

func (self *localProcess) Serialize(ser Serializer) {
	self.a.(Serializable).OnSerialize(ser)
}

func (self *localProcess) notifySystem(data interface{}) {
	self.Tell(&Message{
		Data:      data,
		SourcePID: self.pid,
		TargetPID: self.pid,
	})
}

func (self *localProcess) CreateRPC(waitCallID int64) *util.Future {

	f := util.NewFuture()

	self.mailbox.Hijack(func(rpcBack interface{}) bool {

		rpcMsg := rpcBack.(*Message)
		if rpcMsg.CallID == waitCallID {

			self.mailbox.Hijack(nil)
			f.Done(rpcMsg)
			return true
		}

		return false
	})

	return f
}

func (self *localProcess) PID() *PID {
	return self.pid
}

func (self *localProcess) Tell(data interface{}) {

	if EnableDebug {
		log.Debugf("#notify %s", data.(Context).String())
	}

	self.mailbox.Push(data)
}

func (self *localProcess) Stop() {

	self.notifySystem(&proto.Stop{})
}

func (self *localProcess) OnRecv(data interface{}) {

	ctx := data.(Context)

	if EnableDebug {
		log.Debugf("#recv %s", data.(Context).String())
	}

	self.a.OnRecv(ctx)
}

func newLocalProcess(a Actor, pid *PID) *localProcess {

	self := &localProcess{
		mailbox: mailbox.NewBounded(10),
		a:       a,
		pid:     pid,
	}

	self.mailbox.Start(self)

	self.notifySystem(&proto.Start{})

	return self
}
