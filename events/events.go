// Package events holds event types and according functions, used by the RV framework
package events

import (
	"errors"
	"github.com/TobiasMende/go-run-verify/helpers"
	"time"
)

// InEventType is an enum describing the type of the in event
type InEventType byte

const (
	// The type of the event is unspecified
	UNSPECIFIED InEventType = iota
	// The event was created on sending a message
	SEND
	// The event was created while receiving a message
	RECEIVE
	// The event was created by application internal loggers
	INTERNAL
)

func (t InEventType) String() string {
	switch t {
	case SEND:
		return "S"
	case RECEIVE:
		return "R"
	case INTERNAL:
		return "I"
	default:
		return "UNSPECIFIED"
	}
}

func (t InEventType) UnmarshalBinary(data []byte) (err error) {
	if len(data) == 0 {
		return errors.New("InEventType: data is too short. Requiring 1 byte")
	}
	val := data[0]
	t = InEventType(val)
	return nil
}

// Type InEvent is the internal representation for events occuring in the outer world
type InEvent struct {
	// time at which the event occured
	Created time.Time
	// time at which the event first handled by the logging layer
	Received time.Time
	// type of the event
	Type InEventType
	// payload of the event
	Content interface{}
}

// Type MonitoringEvent is the representation for events created by the monitoring layer
type MonitoringEvent struct {
	// time of event creation
	Created time.Time
	// curent state of the monitor, when creating this event
	CurrentState interface{}
	// current decission of the monitor, when creating this event
	CurrentDecission helpers.MonitorDecission
	// name of the monitor, that created this event
	MonitorName string
	// description of the monitor, that created this event
	MonitorDescription string
}

// Function NewMonitoringEvent is a helper function, creating a new MonitoringEvent with the current time as creation time
func NewMonitoringEvent(state interface{}, decission helpers.MonitorDecission, name string, description string) (evt *MonitoringEvent) {
	return &MonitoringEvent{time.Now(), state, decission, name, description}
}
