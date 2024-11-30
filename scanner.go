package main

import (
	//"context"
    "fmt"
    "net"
    //"sync"
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
