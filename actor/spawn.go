package actor

func SpawnByFunc(name string, a ActorFunc) *PID {

	return SpawnByInstance(name, a)
}

func SpawnByInstance(name string, a Actor) *PID {

	if !inited {
		panic("Call actor.StartSystem first")
	}

	pid := NewLocalPID(name)

	proc := newLocalProcess(a, pid)

	if err := LocalPIDManager.Add(proc); err != nil {
		panic(err)
	}

	pid.proc = proc

	return proc.pid
}
