package sci
//SCI = Serial Communication Interface

import (
	"encoding/binary"
	"bytes"
)

type Behavior struct {
	Angle float64
	Speed float64
	RotationSpeed float64
	SpeedK float64
	IsHeadRelative bool
	EnableHeadSync bool
}

func NewBehavior() Behavior {
	return Behavior{SpeedK: 1}
}

func(this Behavior) Serialize() (bytes.Buffer, error) {
	
	buffer := bytes.Buffer{}
	
	err := binary.Write(&buffer, binary.BigEndian, this)
	
	if err != nil {
		return bytes.Buffer{}, err
	}
	
	return buffer, nil
}

func (this *Behavior) Deserialize(buffer bytes.Buffer) (*Behavior, error) {
	err := binary.Read(&buffer, binary.BigEndian, this)
	
	if err != nil {
		return &Behavior{}, err
	}
	
	return this, nil
}