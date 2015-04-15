# Go Run Verify!

The essence of this project is the demonstration of a runtime verification pipeline, 
i.e. a modular runtime verification system, in Go.

The system currently consists of three layers:

1. *The logging layer*
    The logging layer is responsible for receiving or extracting events from the system under scrutiny.
1. *The monitoring layer*
    The monitoring layer receives events from the logging layer and may have multiple concurrent monitors, that verify system properties independently from each other.
1. *The diagnosis layer*
    The diagnosis layer receives status updates from the monitors and may extract a reason for the monitors verdict.

## Important Notes

This project is in an early stage and just meant to show, how neat the implementation could be, using the Go channels and routines.
Currently, there are no real monitors implemented on the monitoring layer. Furthermore, on the logging layer only exists a logging stub, which generates events. 
In the future, it is planed to develop an interceptor for the CoAP protocol in order to be able to inject the RV system into the CoAP stack.

Furthermore, the framework should be extended by a fourth layer (mitigation layer) later on, which might provide healing-abilities.
