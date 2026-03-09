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
	ViewAngle	float32 `json:"ViewAngle"`
	MotionAngle  float32 `json:"MotionAngle"`
	Speed          float32 `json:"speed"`
	SpeedK         float32 `json:"speedK"`
	EnableHeadSync bool    `json:"enableHeadSync"`
}

func NewBehavior() Behavior {
	return Behavior{SpeedK: 1}
}

func (Behavior) Size() int {
	return 4*4 + 1*1
}

func (this Behavior) Serialize() []byte {

	buffer := make([]byte, this.Size())

	binary.LittleEndian.PutUint32(buffer[0:4], math.Float32bits(this.ViewAngle))
	binary.LittleEndian.PutUint32(buffer[8:12], math.Float32bits(this.MotionAngle))
	binary.LittleEndian.PutUint32(buffer[4:8], math.Float32bits(this.Speed))
	binary.LittleEndian.PutUint32(buffer[12:16], math.Float32bits(this.SpeedK))

	if this.EnableHeadSync {
		buffer[16] = 1
	} else {
		buffer[16] = 0
	}

	return buffer
}

func (this *Behavior) Deserialize(buffer []byte) (*Behavior, error) {
	if len(buffer) < this.Size() {
		return this, TooSmallBufferError
	}

	this.ViewAngle = math.Float32frombits(binary.LittleEndian.Uint32(buffer[0:4]))
	this.MotionAngle = math.Float32frombits(binary.LittleEndian.Uint32(buffer[8:12]))
	this.Speed = math.Float32frombits(binary.LittleEndian.Uint32(buffer[4:8]))
	this.SpeedK = math.Float32frombits(binary.LittleEndian.Uint32(buffer[12:16]))

	this.EnableHeadSync = buffer[16] == 1

	return this, nil
}

func ValidateBehaviorBuffer(buffer []byte) (Behavior, error) {

	behavior, err := (&Behavior{}).Deserialize(buffer)

	if err != nil {
		return Behavior{}, fmt.Errorf("Behavior error: %w: %w", InvalidBehaviorFormat, err)
	}

	return *behavior, nil
}
