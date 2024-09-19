# Charybdis Monitoring System client module

The client module provides a set of functions for transferring metrics and logs to the monitoring system.

## Example

```go
package main

import (
    "log"
    "github.com/DemianSV/chrdsclient"
)

func main() {
    chrdsclient.Conf.SpaceID = "SPACEID"
    chrdsclient.Conf.ModulID = "MODULID"
    chrdsclient.Conf.DataManagerURL = []string{"https://DATAMANAGERURL01", "https://DATAMANAGERURL02"}
    chrdsclient.Conf.DataManagerTimeOut = 1
    chrdsclient.Conf.ClientInSecureSkipVerify = true


    err := chrdsclient.Log("log", "Test test test") // Parameter 1: Metric name, Parameter 2: Value (string).
    if err != nil {
        log.Print(err)
    }

    err := chrdsclient.Metric("test", float32(200)) // Parameter 1: Metric name, Parameter 2: Value (float32).
    if err != nil {
        log.Print(err)
    }
}
```
