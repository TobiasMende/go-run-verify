package stubs

import (
	"github.com/TobiasMende/go-run-verify/events"
	"github.com/TobiasMende/go-run-verify/logging"
	"time"
)

type DemoLogConnector int

func (con DemoLogConnector) Init() error {
	logging.Info.Println("Init DemoLogConnector")
	return nil
}

func (con DemoLogConnector) Receive() (interface{}, error) {
	con++
	time.Sleep(500 * time.Millisecond)
	return con, nil
}

func (con DemoLogConnector) Close() error {
	logging.Info.Println("Closing DemoLogConnector")
	return nil
}

func unmarshalDemo(msg interface{}) (evt events.InEvent, err error) {
	evt.Received = time.Now()
	time.Sleep(100 * time.Millisecond)
	evt.Created = time.Now()
	evt.Type = events.SEND
	evt.Content = msg

	return evt, nil
}

func DemoLogHandler(outChannel chan<- events.InEvent, done <-chan bool) {
	var connector DemoLogConnector
	logging.LogHandler(connector, unmarshalDemo, outChannel, done)
}
