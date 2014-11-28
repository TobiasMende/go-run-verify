package monitoring

import (
	"github.com/TobiasMende/go-run-verify/events"
	"github.com/TobiasMende/go-run-verify/helpers"
	"time"
)

var (
	Trace, Info, Warning, Error = helpers.NewLayerLogger("monitoring")
)

type Monitor func(*MonitorConfiguration, chan<- *MonitorConfiguration)
type MonitorConfiguration struct {
	In      chan *events.InEvent
	Out     chan *events.MonitoringEvent
	Info    MonitorInfo
	Handler Monitor
}

type MonitorInfo struct {
	Name        string
	Description string
}

func NewMonitorConfiguration(monitor Monitor) (config *MonitorConfiguration) {
	Trace.Println("creating new monitor config")
	config = &MonitorConfiguration{}
	config.In = make(chan *events.InEvent, 10)
	config.Out = make(chan *events.MonitoringEvent, 5)
	config.Handler = monitor
	return config
}

func Dispatcher(monitors []*MonitorConfiguration, in <-chan *events.InEvent) {
	terminationChannel := make(chan *MonitorConfiguration, 1)

	// Start the monitors:
	Trace.Print("Starting Monitors...")
	for _, m := range monitors {
		go m.Handler(m, terminationChannel)
	}
	Trace.Println(" done")

End:
	for {
		select {
		case evt, ok := <-in:
			if ok {
				for _, m := range monitors {
					// TODO dispatching must be done in parallel (otherwise, slow monitors would affect the dispatching)
					// alternative: select, default-case: put that channel to the end and report an error if the problem remains
					m.In <- evt

				}
			} else {
				Trace.Print("Dispatcher: in channel closed, delegating...")
				for _, m := range monitors {
					close(m.In)
				}
				Trace.Println(" done")
				break End
			}
		case m := <-terminationChannel:
			Info.Println("receiving termination request from ", m.Info)
			if !deleteMonitor(monitors, m) {
				break End
			}
		}
	}
	Info.Println("Terminating Dispatcher")

}

func deleteMonitor(monitors []*MonitorConfiguration, m *MonitorConfiguration) bool {
	index := -1
	for i, monitor := range monitors {
		if m == monitor {
			index = i
			break
		}
	}
	if index >= 0 {
		// Delete found monitor (by putting it to the end and reducing the size)
		monitors[index], monitors[len(monitors)-1], monitors = monitors[len(monitors)-1], nil, monitors[:len(monitors)-1]
	}

	if len(monitors) == 0 {
		return false
	}
	return true
}

func (config *MonitorConfiguration) PublishEvent(state interface{}, decission helpers.MonitorDecission) {
	config.Out <- &events.MonitoringEvent{time.Now(), state, decission, config.Info.Name, config.Info.Description}
}

func Watchdog(config *MonitorConfiguration, terminationRequest chan<- *MonitorConfiguration) {
	t := 5 * time.Second
Loop:
	for {
		select {
		case _, ok := <-config.In:
			if !ok {
				Trace.Println("Terminating Watchdog on Channel close")
				break Loop

			}
			config.PublishEvent("Everything fine so far", helpers.UN)
		case <-time.After(t):
			config.PublishEvent("Timeout!!!", helpers.BOTTOM)
			break Loop

		}
	}
	terminationRequest <- config
}
