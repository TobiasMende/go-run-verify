package logging

import (
	"github.com/TobiasMende/go-run-verify/events"
	"github.com/TobiasMende/go-run-verify/helpers"
)

var (
	Trace, Info, Warning, Error = helpers.NewLayerLogger("logging")
)

type Unmarshaler func(interface{}) (events.InEvent, error)

type LogConnector interface {
	Init() error
	Receive() (interface{}, error)
	Close() error
}

func LogHandler(connector LogConnector, unmarshal Unmarshaler, outChannel chan<- *events.InEvent, done <-chan bool) {
	err := connector.Init()
	if err != nil {
		Error.Println("Conection failed! Unable to continue: ", err)
		return
	}
Loop:
	for {
		select {
		case <-done:
			break Loop
		default:
			msg, err := connector.Receive()

			if err != nil {
				Error.Println("Failed to receive from LogConnector. ", err)
				break
			}
			evt, err := unmarshal(msg)
			if err != nil {
				Warning.Println("Failed to unmarshal message. Contining: ", err)
				continue // Invalid message, wait for next
			}
			outChannel <- &evt
		}
	}
	connector.Close()
	Info.Println("Terminating LogHandler successful")
}
