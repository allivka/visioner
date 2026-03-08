package sci

import (
	"encoding/binary"
	"fmt"
	"math"
)

var (
	InvalidBehaviorFormat error = fmt.Errorf("Invalid format of behavior package")
	TooSmallBufferError   error = fmt.Errorf("Too small buffer")
)

type Behavior struct {
	Angle          float32 `json:"angle"`
	Speed          float32 `json:"speed"`
	RotationSpeed  float32 `json:"rotationSpeed"`
	SpeedK         float32 `json:"speedK"`
	IsHeadRelative bool    `json:"isHeadRelative"`
	EnableHeadSync bool    `json:"enableHeadSync"`
}

func NewBehavior() Behavior {
	return Behavior{SpeedK: 1}
}

func (Behavior) Size() int {
	return 4*4 + 1*2
}

func (this Behavior) Serialize() []byte {

	buffer := make([]byte, this.Size())

	binary.LittleEndian.PutUint32(buffer[0:4], math.Float32bits(this.Angle))
	binary.LittleEndian.PutUint32(buffer[4:8], math.Float32bits(this.Speed))
	binary.LittleEndian.PutUint32(buffer[8:12], math.Float32bits(this.RotationSpeed))
	binary.LittleEndian.PutUint32(buffer[12:16], math.Float32bits(this.SpeedK))

	if this.IsHeadRelative {
		buffer[16] = 1
	} else {
		buffer[17] = 0
	}
	

	if this.EnableHeadSync {
		buffer[16] = 1
	} else {
		buffer[17] = 0
	}

	return buffer
}

func (this *Behavior) Deserialize(buffer []byte) (*Behavior, error) {
	if len(buffer) < this.Size() {
		return this, TooSmallBufferError
	}

	this.Angle = math.Float32frombits(binary.LittleEndian.Uint32(buffer[0:4]))
	this.Speed = math.Float32frombits(binary.LittleEndian.Uint32(buffer[4:8]))
	this.RotationSpeed = math.Float32frombits(binary.LittleEndian.Uint32(buffer[8:12]))
	this.SpeedK = math.Float32frombits(binary.LittleEndian.Uint32(buffer[12:16]))

	this.IsHeadRelative = buffer[16] == 1
	this.EnableHeadSync = buffer[17] == 1

	return this, nil
}

func ValidateBehaviorBuffer(buffer []byte) (Behavior, error) {

	behavior, err := (&Behavior{}).Deserialize(buffer)

	if err != nil {
		return Behavior{}, fmt.Errorf("Behavior error: %w: %w", InvalidBehaviorFormat, err)
	}

	return *behavior, nil
}
