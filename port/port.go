package port

import (
    "fmt"
    "net"
    "time"
    "strings"
    "unicode"
    "strconv"
    "sort"
)

type ScanResult struct {
    Port int
    Protocol string
    State string
}

type ScanResults []ScanResult

func (sr ScanResults) Len() int {
    return len(sr)
}

func (sr ScanResults) Less(i, j int) bool {
    return sr[i].Port < sr[j].Port
}

func (sr ScanResults) Swap(i, j int) {
    sr[i], sr[j] = sr[j], sr[i]
}

func ScanPort(protocol, hostname string, port int) ScanResult {
    result := ScanResult{Port:port,Protocol:protocol}
    address := fmt.Sprintf("%s:%d", hostname, port)
    conn, err := net.DialTimeout(protocol, address, 60*time.Second)
    // TODO: add logic to determine closed or filtered
    if err != nil {
        if strings.Contains(fmt.Sprint(err), "connection refused") {
            result.State = "Closed"
        } else {
            result.State = "Filtered"
        }
        return result
    }
    defer conn.Close()
    result.State = "Open"
    return result
}

func Scanner(protocol, hostname string, port_c chan int, result_c chan ScanResult) {
    for port := range port_c {
        result_c <- ScanPort(protocol, hostname, port)
    }
}

func Scan(hostname, protocol string, prange string) []ScanResult {
    var results ScanResults

    f := func(c rune) bool {
        return !unicode.IsNumber(c)
    }
    pr := strings.FieldsFunc(prange, f)
    if len(pr) != 2 {
        fmt.Println("Range format error.")
        return results
    }
    startp, _ := strconv.Atoi(string(pr[0]))
    endp, _ := strconv.Atoi(string(pr[1]))
    if endp > 65535 {
        endp = 65535
        fmt.Println("There are only a total of 65535 ports.")
    }

    buffer := 200
    port_c := make(chan int, buffer)
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

    sort.Sort(results)
    return results
}