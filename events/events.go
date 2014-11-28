package events

import (
	"encoding/binary"
	"errors"
	"github.com/TobiasMende/go-run-verify/helpers"
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
	CurrentState       interface{} // TODO chose concrete type for state
	CurrentDecission   helpers.MonitorDecission
	MonitorName        string
	MonitorDescription string
}

func NewMonitoringEvent(state interface{}, decission helpers.MonitorDecission, name string, description string) (evt *MonitoringEvent) {
	return &MonitoringEvent{time.Now(), state, decission, name, description}
}
