package main

import (
	"flag"
	"fmt"
	"regexp"
)

func main() {
	ptcl := flag.String("type", "tcp", "the protocol to use ")
	host := flag.String("host", "127.0.0.1", "The target host to scan")
	ports := flag.String("ports", "20-1000", "list of ports to scan ex : 80,443 or 0-100")
	timeout := flag.Int("timeout", 5, "Timeout in seconds for the connection")
	flag.Parse()

	portRange, err := parsePorts(*ports)

	var singleHost string
	var multiHost []string 
	is_singleHost := false 

	if checkIP(*host) {
		finalHost = *host
		is_singleHost := true
	} else {
		if matched, _ := regexp.MatchString(`^([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`, *host); matched {
			ip ,err:= reverselookup(*host)
			if err != nil {
				fmt.Printf("Erreur de reverse lookup pour l'URL %s : %v\n", *host, err)
				return
			}
			multiHost = ip
		} else {
			fmt.Println("Erreur : L'adresse ou l'URL fournie n'est pas valide.")
			return
		}
	}

	if !checkProtocole(*ptcl) {
		fmt.Println("Error with the protocol")
		return
	}
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if is_singleHost{
		scanAndGrab(finalHost, portRange, *timeout)
	} else {
		scanIPRange(multiHost, portRange,*timeout)
	}

	
}