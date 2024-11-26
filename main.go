package main

// 2. Importation des packages nécessaires
import (
	"flag"
	"fmt"
)

func main() {
	// Définir les options
	ptcl := flag.String("type","tcp","the protocol to use ")
	host := flag.String("host", "127.0.0.1", "The target host to scan")
	ports := flag.String("ports", "0-1000", "list of ports to scan ex : 80,443 or 0-100")
	timeout := flag.Int("timeout", 5, "Timeout in seconds for the connection")



	port_range, err := parsePorts(*ports)
	checkAddr = checkIP(host)
	checkPtcl = checkProtocole(ptcl) 

	if !checkAddr{
		fmt.Println("Error with the host address")
		return
	}
	if !checkPtcl{
		fmt.Println("Error with the protocol")
		return 
	}
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	result := scanner("tcp", *host, port_range, *timeout)
	


}
