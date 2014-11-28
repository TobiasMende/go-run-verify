// Package monitoring contains the RV framework architecture concerning the monitoring layer.
//
// RV Framework Layers:  logging -> MONITORING -> diagnosis -> mitigation
package monitoring

import (
	"github.com/TobiasMende/go-run-verify/events"
	"github.com/TobiasMende/go-run-verify/helpers"
	"time"
)

var (
	Trace, Info, Warning, Error = helpers.NewLayerLogger("monitoring")
)

// Type Monitor is the type of functions that should be executed as runtime verification monitors
type Monitor func(*MonitorConfiguration, chan<- *MonitorConfiguration)

// Type MonitorConfiguration holds the configuration of a single monitor
type MonitorConfiguration struct {
	// channel for receiving events from the logging layer
	In chan *events.InEvent
	// channel for sending events to the diagnosis layer
	Out chan *events.MonitoringEvent
	// information about the monitor (human readable)
	Info MonitorInfo
	// the monitor function itself
	Handler Monitor
}

// Type MonitorInfo contains information about the monitor.
// The information is supposed to be human readable. It is useful, when logging information about the monitor.
type MonitorInfo struct {
	// The name of the monitor. E.g. "SendReceiveOrderMonitor"
	Name string
	// The description of the monitor, e.g. a LTL formula or a longer text, describing the sense of the monitor
	Description string
}

// Function NewMonitorConfiguration is a helper function for creating a new MonitorConfiguration with pre-initialized channels and a given Monitor function
func NewMonitorConfiguration(monitor Monitor) (config *MonitorConfiguration) {
	Trace.Println("creating new monitor config")
	config = &MonitorConfiguration{}
	config.In = make(chan *events.InEvent, 10)
	config.Out = make(chan *events.MonitoringEvent, 5)
	config.Handler = monitor
	return config
}

// Dispatcher function is the entry point for events comming from the logging layer
// The dispatcher should be executed as goroutine. Its job is the delegation of incoming events to alle monitors.
// The dispatcher first executes all monitors as goroutines and than waits for incoming events for dispatching
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

// Helper function deleteMonitor deletes a monitor m from the array monitors.
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

// Method PublishEvents sends an event to the diagnosis layer
func (config *MonitorConfiguration) PublishEvent(state interface{}, decission helpers.MonitorDecission) {
	config.Out <- &events.MonitoringEvent{time.Now(), state, decission, config.Info.Name, config.Info.Description}
}

// Watchdog is an example for a monitor.
// The watchdog waits for incoming events and send the Decission BOTTOM, when no events are incoming for more than 5 seconds.
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
