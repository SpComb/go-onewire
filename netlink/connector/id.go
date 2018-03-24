package connector

type ID struct {
	Idx uint32
	Val uint32
}

func (id *ID) MarshalBinary() ([]byte, error) {
	return marshalBinary(*id)
}

func (id *ID) UnmarshalBinary(data []byte) error {
	return unmarshalBinary(id, data)
}
