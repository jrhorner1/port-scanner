package port

import (
    "fmt"
    "net"
    "time"
    "strings"
    "strconv"
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

func Scan(hostname, protocol string, prange string) []ScanResult {
    var results []ScanResult

    pr := strings.Split(prange, "-")
    if len(prange) == 1 {
        fmt.Println("Range must be specified as `1-10`.")
        // exit
    }
    startp, _ := strconv.Atoi(string(pr[0]))
    endp, _ := strconv.Atoi(string(pr[1]))

    port_c := make(chan int, 100)
    result_c := make(chan ScanResult)

    for i := 0; i < cap(port_c)/2; i++ {
        go Scanner(protocol, hostname, port_c, result_c)
    }

    go func() {
        for i := startp; i <= endp; i++ {
            port_c <- i
        }
    }()

    for i := startp; i <= endp; i++ {
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