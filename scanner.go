package main

import (
	"context"
    "fmt"
    "net"
    "sync"
    "time"
)



func grabBanner(conn net.Conn, port int) string {
    // Utilisez la nouvelle fonction advancedBannerGrab
    return advancedBannerGrab(conn, port)
}

// Structure pour stocker les résultats du scan
type ScanResult struct {
    Port    int
    State   string
    Banner  string
}

/*func parallelScan(ip string, ports []int, timeout int, grab bool) map[string]string {
    results := make(map[string]string)
    resultsChan := make(chan ScanResult, len(ports))
    errChan := make(chan error, 1)
    var wg sync.WaitGroup
    var mu sync.Mutex
    semaphore := make(chan struct{}, 100)

    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
    defer cancel()

    sentPorts := make(map[int]bool)

    for _, port := range ports {
        wg.Add(1)
        semaphore <- struct{}{}

        go func(port int) {
            defer wg.Done()
            defer func() { <-semaphore }()
            conn, err := (&net.Dialer{
                Timeout: time.Duration(timeout) * time.Second,
            }).DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", ip, port))
            
            if err != nil {
                return
            }
            defer conn.Close()

            result := ScanResult{
                Port:   port,
                State:  "open",
                Banner: "",
            }
            if grab {
                result.Banner = grabBanner(conn, port)
            }
            select {
            case resultsChan <- result:
            case <-ctx.Done():
                return
            }
        }(port)
    }

    go func() {
        wg.Wait()
        close(resultsChan)
        close(errChan)
    }()

    for result := range resultsChan {
        mu.Lock()
        if !sentPorts[result.Port] {
            protPort := fmt.Sprintf("%d/tcp", result.Port)
            if result.Banner != "" {
                results[protPort] = fmt.Sprintf("open\n%s", result.Banner)
            } else {
                results[protPort] = "open"
            }
            sentPorts[result.Port] = true
        }
        mu.Unlock()
    }
    return results
}*/


func scan(ip string, ports []int, timeout int, grab bool) map[string]string {
    toReturn := make(map[string]string)
    
    for _, port := range ports {
        address := fmt.Sprintf("%s:%d", ip, port)
        conn, err := net.DialTimeout("tcp", address, time.Duration(timeout)*time.Second)
        if err != nil {
            continue // Port fermé ou non accessible
        }
        defer conn.Close()

        protPort := fmt.Sprintf("%d/tcp", port)
        
        if grab {
            banner := grabBanner(conn, port)
            toReturn[protPort] = fmt.Sprintf("open\n%s", banner)
        } else {
            toReturn[protPort] = "open"
        }
    }
    
    return toReturn
}

// Convertit une map en une chaîne lisible.
func toString(input map[string]string) {
	for key, value := range input {
		fmt.Printf("%s : %s \n", key,value )
	}
}

func scanIPRange(ips []string , ports []int , timeout int,grab bool){
	for _,ip := range ips {
		fmt.Printf("%s : \n", ip)
		toString(scan(ip,ports,timeout,grab))
	}
}
