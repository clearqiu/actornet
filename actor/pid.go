package actor

type PID struct {
	Address string
	Id      string

	proc Process
}

func (self *PID) IsLocal() bool {
	return LocalPIDManager.Address == self.Address
}

func (self *PID) raw() PID {

	return PID{
		Address: self.Address,
		Id:      self.Id,
	}
}

func (self *PID) ref() Process {

	if self.proc != nil {
		return self.proc
	}

	if self.IsLocal() {
		// 更新Process缓冲
		p := LocalPIDManager.Get(self.Id)
		if p != nil {
			self.proc = p
			return p
		}

	} else if RemoteProcessCreator != nil {
		self.proc = RemoteProcessCreator()
		return self.proc
	}

	return nil
}

func (self *PID) Send(target *PID, data interface{}) {

	if target != nil {
		target.ref().Send(self, data)
	} else {
		panic("empty target")
	}

}

func (self *PID) String() string {
	if self == nil {
		return "nil"
	}
	return self.Address + "/" + self.Id
}

func NewPID(address, id string) *PID {
	return &PID{
		Address: address,
		Id:      id,
	}
}

func NewLocalPID(id string) *PID {
	return &PID{
		Address: LocalPIDManager.Address,
		Id:      id,
	}
}

var Root = NewLocalPID("Root")

var RemoteProcessCreator func() Process
