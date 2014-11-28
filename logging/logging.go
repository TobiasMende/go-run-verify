// Package logging provides the logging layer, speaking in runtime verification terms.
// logging is the first layer, which receives information about the application state from the outer world
//
// RV Framework Layers: LOGGING -> monitoring -> diagnosis -> mitigation

package logging

import (
	"github.com/TobiasMende/go-run-verify/events"
	"github.com/TobiasMende/go-run-verify/helpers"
)

var (
	Trace, Info, Warning, Error = helpers.NewLayerLogger("logging")
)

// Type Unmarshaler is the type, of which functions must be, that the LogHandler uses to generate an InEvent from a marshaled representation
type Unmarshaler func(interface{}) (events.InEvent, error)

// Type LogConnecter is the interface, which should be implemented by types that are able to create a connection to the outer world for receiving marshaled InEvents.
type LogConnector interface {
	// Function Init is called exactly once, when starting the LogHandler
	Init() error

	// Function Receive is called whenever the LogHandler is ready for receiving. This function should block, if no incoming data is available
	Receive() (interface{}, error)

	// Function Close is called exactly once, when the LogHandler is terminating
	Close() error
}

//Function LogHandler handles the process of receiving and unmarshaling data and sends InEvents to the RV framework
//LogHandler should be executed as goroutine and can be terminated by sending a value via the done channel
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
			close(outChannel)
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
