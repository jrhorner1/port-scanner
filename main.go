package main

import (
    "github.com/jrhorner1/port-scanner/port"
    "fmt"
)

func main() {
    scan := port.InitialScan("192.168.11.2", "tcp")
    for i := range scan {
        fmt.Printf("%d/%s %s\n",scan[i].Port, scan[i].Protocol, scan[i].State)
    }
}