package sci

import (
	"fmt"
	"unsafe"
)

var (
	InvalidBehaviorFormat error = fmt.Errorf("Invalid format of behavior package")
	TooSmallBufferError error = fmt.Errorf("Too small buffer")
)

type Behavior struct {
	Angle          float64 `json:"angle"`
	Speed          float64 `json:"speed"`
	RotationSpeed  float64 `json:"rotationSpeed"`
	SpeedK         float64 `json:"speedK"`
	IsHeadRelative bool    `json:"isHeadRelative"`
	EnableHeadSync bool    `json:"enableHeadSync"`
}

func NewBehavior() Behavior {
	return Behavior{SpeedK: 1}
}

func (Behavior) Size() int {
	return 8 * 4 + 1 * 2;
}

func (this Behavior) Serialize() ([]byte) {
	
	buffer := make([]byte, this.Size())
	
	ptr := unsafe.SliceData(buffer)
	
	f64ptr := (*float64)(unsafe.Pointer(ptr))
	
	*f64ptr = this.Angle
	f64ptr = (*float64)(unsafe.Add(unsafe.Pointer(f64ptr), 8))
	
	*f64ptr = this.Speed
	f64ptr = (*float64)(unsafe.Add(unsafe.Pointer(f64ptr), 8))
	
	*f64ptr = this.RotationSpeed
	f64ptr = (*float64)(unsafe.Add(unsafe.Pointer(f64ptr), 8))
	
	*f64ptr = this.SpeedK
	f64ptr = (*float64)(unsafe.Add(unsafe.Pointer(f64ptr), 8))
	
	boolPtr := (*bool)(unsafe.Pointer(f64ptr))
	
	*boolPtr = this.IsHeadRelative
	boolPtr = (*bool)(unsafe.Add(unsafe.Pointer(boolPtr), 1))
	
	*boolPtr = this.EnableHeadSync
	
	
	return buffer
}

func (this *Behavior) Deserialize(buffer []byte) (*Behavior, error) {
	if len(buffer) < this.Size() {
		return this, TooSmallBufferError
	}
	
	ptr := unsafe.SliceData(buffer)
	
	f64ptr := (*float64)(unsafe.Pointer(ptr))
	
	this.Angle = *f64ptr
	f64ptr = (*float64)(unsafe.Add(unsafe.Pointer(f64ptr), 8))
	
	this.Speed = *f64ptr
	f64ptr = (*float64)(unsafe.Add(unsafe.Pointer(f64ptr), 8))
	
	this.RotationSpeed = *f64ptr
	f64ptr = (*float64)(unsafe.Add(unsafe.Pointer(f64ptr), 8))
	
	this.SpeedK = *f64ptr
	f64ptr = (*float64)(unsafe.Add(unsafe.Pointer(f64ptr), 8))
	
	boolPtr := (*bool)(unsafe.Pointer(f64ptr))
	
	this.IsHeadRelative = *boolPtr
	boolPtr = (*bool)(unsafe.Add(unsafe.Pointer(boolPtr), 1))
	
	this.EnableHeadSync = *boolPtr
	
	return this, nil
}

func ValidateBehaviorBuffer(buffer []byte) (Behavior, error) {

	behavior, err := (&Behavior{}).Deserialize(buffer)

	if err != nil {
		return Behavior{}, fmt.Errorf("Behavior error: %w: %w", InvalidBehaviorFormat, err)
	}

	return *behavior, nil
}
