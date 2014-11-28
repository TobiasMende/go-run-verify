package diagnosis

import (
	"github.com/TobiasMende/go-run-verify/events"
	"github.com/TobiasMende/go-run-verify/helpers"
)

var (
	Trace, Info, Warning, Error = helpers.NewLayerLogger("diagnosis")
)

func Logger(ec <-chan *events.MonitoringEvent) {
	for {
		evt, more := <-ec
		if more {
			Info.Println(evt)
		} else {
			break
		}
	}
	Info.Println("Logger finished")
}
