package logging

import (
	"encoding/binary"
	"errors"
	"github.com/TobiasMende/go-run-verify/events"
	_ "github.com/dustin/go-coap"
	"time"
)

//UnmarshalCoapMessage expects a byte array as parameter and returns an InEvent
// The expected message format is |1 byte EventType | 8 byte int64 timestamp | CoAP message |
func UnmarshalCoapMessage(msg interface{}) (evt events.InEvent, err error) {
	bytes := msg.([]byte)

	if len(bytes) < 10 {
		return evt, errors.New("UnmarshalCoapMessage: msg is too short")
	}
	var eventType events.InEventType
	eventType.UnmarshalBinary(bytes[:1])
	evt.Type = eventType
	evt.Received = time.Now()

	timestamp, _ := binary.Varint(bytes[1:9])
	evt.Created = time.Unix(timestamp, 0)

	// TODO parse CoAP message
	return evt, nil

}
