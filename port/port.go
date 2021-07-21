package port

import (
    "fmt"
    "net"
    "time"
)

type ScanResult struct {
    Port int
    Protocol string
    State string
}

func ScanPort(protocol, hostname string, port int) ScanResult {
    result := ScanResult{Port:port,Protocol:protocol}
    address := fmt.Sprintf("%s:%d", hostname, port)
    conn, err := net.DialTimeout(protocol, address, 60*time.Second)
    if err != nil {
        result.State = "Closed"
        return result
    }
    defer conn.Close()
    result.State = "Open"
    return result
}

func InitialScan(hostname, protocol string) []ScanResult {
    var results []ScanResult

    port_c := make(chan int, 100)
    result_c := make(chan ScanResult)

    for i := 0; i < cap(port_c)/2; i++ {
        go Scanner(protocol, hostname, port_c, result_c)
    }

    go func() {
        for i := 1; i <= 1024; i++ {
            port_c <- i
        }
    }()

    for i := 0; i < 1024; i++ {
        result := <- result_c
        results = append(results, result)
    }

    close(port_c)
    close(result_c)

    return results
}

func Scanner(protocol, hostname string, port_c chan int, result_c chan ScanResult) {
    for port := range port_c {
        result_c <- ScanPort(protocol, hostname, port)
    }
}