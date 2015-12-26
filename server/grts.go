package main

import (
    "log"
    "os"
    "io"
    "io/ioutil"
)

/* Logger */
var (
    Trace   *log.Logger
    Info    *log.Logger
    Warning *log.Logger
    Error   *log.Logger
)

func configureLogger(traceHandle io.Writer, infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer) {
    Trace = log.New(traceHandle, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
    Info = log.New(infoHandle, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
    Warning = log.New(warningHandle, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
    Error = log.New(errorHandle, "ERROR: ",    log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
    configureLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
    server := Server{}
    server.Start()
}
