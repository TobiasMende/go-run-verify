// Package diagnosis holds the implementation of the diagnosis layer
// diagnosis is just ment as stub in the current implementation, hence no diagnosis is done here (just logging)
//
// RV Framework Layers: logging -> monitoring -> DIAGNOSIS -> mitigation
package diagnosis

import (
	"github.com/TobiasMende/go-run-verify/events"
	"github.com/TobiasMende/go-run-verify/helpers"
)

var (
	Trace, Info, Warning, Error = helpers.NewLayerLogger("diagnosis")
)

//Logger waits for incomming MonitoringEvents and logs them until the channel is closed.
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
