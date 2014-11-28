package events

import (
	"encoding/binary"
	"errors"
	"time"
)

type InEventType uint

const (
	SEND InEventType = iota
	RECEIVE
)

func (t InEventType) UnmarshalBinary(data []byte) (err error) {
	if len(data) == 0 {
		return errors.New("InEventType: data is too short. Requiring 1 byte")
	}
	val, _ := binary.Uvarint(data)
	t = InEventType(val)
	return nil
}

type InEvent struct {
	Created  time.Time
	Received time.Time
	Type     InEventType
	Content  interface{}
}

type MonitoringEvent struct {
	Created            time.Time
	CurrentState       interface{}
	CurrentDecission   interface{}
	MonitorName        string
	MonitorDescription string
}
