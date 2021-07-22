package main

import (
    "github.com/jrhorner1/port-scanner/port"
    "fmt"
    "flag"
    "github.com/go-ping/ping"
    "time"
)

func main() {
    var hostname, protocol, portrange string
    flag.StringVar(&hostname, "h", "localhost", "Hostname|IP: example.com | 127.0.0.1")
    flag.StringVar(&protocol, "p", "tcp", "Protocol: tcp | udp")
    flag.StringVar(&portrange, "r", "1-1024", "Port Range: 1-1024")
    flag.Parse()

    // Check if the host is up
    pinger, err := ping.NewPinger(hostname)
    if err != nil {
        panic(err)
    }
    pinger.Count = 3
    pinger.Timeout = 10 * time.Second
    err = pinger.Run()
    if err != nil {
        panic(err)
    }
    stats := pinger.Statistics()
    if stats.PacketLoss == 100 {
        fmt.Println("Host seems to be down.")
        return
    }

    // Perform the scan
    scan := port.Scan(hostname, protocol, portrange)
    for i := range scan {
        if scan[i].State == "Open" || scan[i].State == "Filtered" {
            fmt.Printf("%d/%s %s\n",scan[i].Port, scan[i].Protocol, scan[i].State)
        }
    }
}