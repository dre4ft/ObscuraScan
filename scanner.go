package main

import (
	//"context"
    "fmt"
    "net"
    "sync"
    "time"
)



func grabBanner(conn net.Conn, port int) string {
    return advancedBannerGrab(conn, port)
}




type ScanResult struct {
	Port   int
	State  string
	Banner string
}


func scan(ip string, ports []int, timeout int, grab bool) map[string]string {
	toReturn := make(map[string]string)
	var wg sync.WaitGroup 

	
	var mu sync.Mutex

	for _, port := range ports {
		wg.Add(1) 
		go func(port int) {
			defer wg.Done() 

			address := fmt.Sprintf("%s:%d", ip, port)
			conn, err := net.DialTimeout("tcp", address, time.Duration(timeout)*time.Second)
			if err != nil {
				
				return
			}
			defer conn.Close()

			protPort := fmt.Sprintf("%d/tcp", port)

			
			var result string
			if grab {
				banner := grabBanner(conn, port)
				result = fmt.Sprintf("open\n%s", banner)
			} else {
				result = "open"
			}
			mu.Lock()
			toReturn[protPort] = result
			mu.Unlock()
		}(port)
	}


	wg.Wait()

	return toReturn
}



func toString(input map[string]string) {
	for key, value := range input {
		cve := SearchAndFormat(value)
		fmt.Printf("%s : %s \n", key,value)
		fmt.Printf("Vulnerability scan : \n %s",cve)
		fmt.Printf("-------------------------------")
	}
}

func scanIPRange(ips []string , ports []int , timeout int,grab bool){
	for _,ip := range ips {
		fmt.Printf("%s : \n", ip)
		toString(scan(ip,ports,timeout,grab))
	}
}
