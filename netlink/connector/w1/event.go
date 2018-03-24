package w1

type Event struct {
	Type MessageType
	ID   MessageID
}

func (event Event) MasterID() MasterID {
	var masterID MasterID

	masterID.Unpack(event.ID)

	return masterID
}

func (event Event) SlaveID() SlaveID {
	var slaveID SlaveID

	slaveID.Unpack(event.ID)

	return slaveID
}
