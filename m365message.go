// Package m365message provides methods to build messages
// that can be sent to a Xiaomi M365 electric scooter via BLE
package m365message

import (
	"errors"
	"fmt"
)

// ReadWrite values determine if data is
// to be written to or read from the scooter
type ReadWrite int

// Direction values determine the flow of data
// on the scooter
type Direction int

const (
	// READ from scooter
	READ ReadWrite = 0x01
	// WRITE to scooter
	WRITE ReadWrite = 0x03
	// MASTER_TO_M365 communicate with M365
	MASTER_TO_M365 Direction = 0x20
	// MASTER_TO_BATTERY communicate with battery
	MASTER_TO_BATTERY Direction = 0x22
)

var (
	// ErrorInvalidDirection direction must be MASTER_TO_M365 or MASTER_TO_BATTERY
	ErrorInvalidDirection = errors.New("invalid direction")
	// ErrorInvalidReadWrite must be READ or WRITE
	ErrorInvalidReadWrite = errors.New("invalid readwrite")
	// ErrorPayloadRequired a payload is required
	ErrorPayloadRequired = errors.New("payload required")
	// ErrorPositionRequired a position is required
	ErrorPositionRequired = errors.New("position required")
)

// Message contains all the fields required to Build() a Message
type Message struct {
	rw        ReadWrite
	payload   []int
	position  int
	direction Direction
}

// NewMessage returns a pointer to a new Message
func NewMessage() *Message {
	return &Message{}
}

// SetRW modifies the rw field on the Message
func (m *Message) SetRW(i ReadWrite) {
	m.rw = i
}

// SetPayload sets the payload field on the Message
// payload is a slice of ints.
// Each call replaces any previously set payload
func (m *Message) SetPayload(payload []int) {
	m.payload = nil
	m.payload = append(m.payload, payload...)
}

// SetPosition sets the position field on the Message
func (m *Message) SetPosition(i int) {
	m.position = i
}

// SetDirection sets the direction field on the Message
func (m *Message) SetDirection(i Direction) {
	m.direction = i
}

// generateChecksum returns a checksum per the spec as an int
func (m *Message) generateChecksum() int {
	var checksum int
	checksum += int(m.direction)
	checksum += m.position
	checksum += int(m.rw)
	checksum += len(m.payload) + 2
	for _, v := range m.payload {
		checksum += v
	}
	return checksum
}

// Build generates a string representation of the Message for
// transmission via BLE to the electric scooter
func (m *Message) Build() (string, error) {
	if m.rw != READ && m.rw != WRITE {
		return "", ErrorInvalidReadWrite
	}

	if m.direction != MASTER_TO_BATTERY && m.direction != MASTER_TO_M365 {
		return "", ErrorInvalidDirection
	}

	if m.payload == nil {
		return "", ErrorPayloadRequired
	}

	if m.position == 0 {
		return "", ErrorPositionRequired
	}

	var tmp []int

	// headers
	tmp = append(tmp, 0x55)
	tmp = append(tmp, 0xAA)

	// body
	tmp = append(tmp, len(m.payload)+2)
	tmp = append(tmp, int(m.direction))
	tmp = append(tmp, int(m.rw))
	tmp = append(tmp, m.position)
	tmp = append(tmp, m.payload...)

	// checksum
	checksum := m.generateChecksum()
	checksum ^= 0xffff // xor
	tmp = append(tmp, checksum&0xff)
	tmp = append(tmp, checksum>>8)

	var res string
	for _, v := range tmp {
		res += fmt.Sprintf("%02x", v)
	}

	return res, nil
}
