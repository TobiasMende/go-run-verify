package helpers

import (
	"fmt"
	"log"
	"os"
)

func NewLayerLogger(layer string) (trace *log.Logger, info *log.Logger, warn *log.Logger, err *log.Logger) {
	trace = log.New(os.Stdout, fmt.Sprint("TRACE ", layer, "\t"), log.LstdFlags)
	info = log.New(os.Stdout, fmt.Sprint("INFO  ", layer, "\t"), log.LstdFlags)
	warn = log.New(os.Stdout, fmt.Sprint("WARN  ", layer, "\t"), log.LstdFlags)
	err = log.New(os.Stderr, fmt.Sprint("ERROR ", layer, "\t"), log.LstdFlags)
	return
}
