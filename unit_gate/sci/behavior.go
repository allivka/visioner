package sci

import (
	"bytes"
	"encoding/binary"
)

type Behavior struct {
	Angle float64 `json:"angle"`
	Speed float64 `json:"speed"`
	RotationSpeed float64 `json:"rotationSpeed"`
	SpeedK float64 `json:"speedK"`
	IsHeadRelative bool `json:"isHeadRelative"`
	EnableHeadSync bool `json:"enableHeadSync"`
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