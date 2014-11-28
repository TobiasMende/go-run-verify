package main

import (
	"fmt"
	"github.com/TobiasMende/go-run-verify/diagnosis"
	"github.com/TobiasMende/go-run-verify/events"
	"github.com/TobiasMende/go-run-verify/stubs"
	"time"
)

func main() {
	loggingDemo()
}

func loggingDemo() {
	ec := make(chan events.InEvent, 1)
	done := make(chan bool, 1)

	go stubs.DemoLogHandler(ec, done)
	for i := 0; i < 10; i++ {
		evt := <-ec
		fmt.Println(i, evt)
	}
	done <- true
	time.Sleep(1 * time.Second)
}

func diagnosisDemo() {
	c := make(chan events.MonitoringEvent, 1)
	go diagnosis.Logger(c)

	for i := 0; i < 10; i++ {
		var evt events.MonitoringEvent
		evt.Created = time.Now()
		evt.MonitorName = fmt.Sprint("Monitor", i)
		c <- evt
		time.Sleep(1 * time.Second)
	}
	close(c)
	time.Sleep(1 * time.Second)
}
