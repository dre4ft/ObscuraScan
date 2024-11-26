
package main

// 2. Importation des packages nécessaires
import (
	"net"
    "fmt"
	"time"
)
   

func whosup( address string, portRange []int, timeout int ) map[string]bool {
    toReturn := make(map[string]bool)
    duration := time.Duration(timeout) * time.Second

    for _, port := range portRange {       
        toReturn[protocol_port] = scan("tcp",address,port)
    }
    return toReturn
}

func scan(protocol string,ip string,port int , timeout int)bool{
	ip_port := fmt.Sprintf("%s:%d", address, port)
        protocol_port := fmt.Sprintf("%s,%d", protocol, port)
        conn, err := net.DialTimeout(protocol, ip_port, duration)
        if err != nil {
			conn.Close()
            return false
        }
		conn.Close()
		return true 
}









func toString(input map[string]bool) {
	for key, value := range input { 
		if value {
			fmt.Printf("%s : open\n", key) 
		} else {
			fmt.Printf("%s : closed\n", key) 
		}
	}
}

