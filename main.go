package main

import (
    "github.com/jrhorner1/port-scanner/port"
    "fmt"
)

func main() {
    scan := port.Scan("scanme.nmap.org", "tcp", "50-250")
    for i := range scan {
        if scan[i].State == "Open" {
            fmt.Printf("%d/%s %s\n",scan[i].Port, scan[i].Protocol, scan[i].State)
        }
    }
}