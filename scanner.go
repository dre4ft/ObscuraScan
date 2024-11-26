
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
		protocol_port := fmt.Sprintf("tcp,%d",port) 
		is_up :=   scan("tcp",address,port,duration)
		if is_up{
			toReturn[protocol_port] = is_up
		} 
       
    }
    return toReturn
}

func scan(protocol string,ip string,port int , timeout time.Duration)bool{
	ip_port := fmt.Sprintf("%s:%d", ip, port)
    conn, err := net.DialTimeout(protocol, ip_port, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true 

}



/*// Lire la bannière
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	reader := bufio.NewReader(conn)
	banner, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("Erreur de lecture : %v", err)
	}

	return banner*/








func toString(input map[string]bool) {
	for key, value := range input { 
		if value {
			fmt.Printf("%s : open\n", key) 
		} else {
			fmt.Printf("%s : closed\n", key) 
		}
	}
}

