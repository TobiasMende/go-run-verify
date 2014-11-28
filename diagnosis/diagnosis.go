package diagnosis

import (
	"github.com/TobiasMende/go-run-verify/events"
	"github.com/TobiasMende/go-run-verify/helpers"
)

var (
	Trace, Info, Warning, Error = helpers.NewLayerLogger("diagnosis")
)

func Logger(name string, ec <-chan *events.MonitoringEvent) {
	for {
		evt, more := <-ec
		if more {
			if evt.CurrentDecission == helpers.BOTTOM {
				Warning.Println(name, evt)
			} else {
				Info.Println(name, evt)
			}
		} else {
			break
		}
	}
	Info.Println("Logger finished")
}
