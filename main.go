package main

import (
    "github.com/jrhorner1/port-scanner/port"
    "fmt"
)

func main() {
    scan := port.InitialScan("scanme.nmap.org")
    for s := range scan {
        fmt.Println(s.Port, s.State)
    }
}