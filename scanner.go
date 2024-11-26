
package main

// 2. Importation des packages nécessaires
import (
	"net"
    "fmt"
	"time"
)
   

func scanner(protocol string, address string, portRange []int, timeout int) map[string]bool {
    toReturn := make(map[string]bool)
    duration := time.Duration(timeout) * time.Second

    for _, port := range portRange {
        ip_port := fmt.Sprintf("%s:%d", address, port)
        protocol_port := fmt.Sprintf("%s,%d", protocol, port)
        
    
        conn, err := net.DialTimeout(protocol, ip_port, duration)
        
        if err != nil {
           
            toReturn[protocol_port] = false
            continue
        }
     
        conn.Close()
        
       
        toReturn[protocol_port] = true
    }

    return toReturn
}
func toString(input map[string]bool){
	for key,value in range input{
		if value{
			fmt.Println("‰s : open",value)
		}
		else {
			fmt.Println("‰s : open",value)
		}
		
	}
}

