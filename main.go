package main

import (
	"flag"
	"fmt"
	"net"
	"regexp"
	"time"
)

func main() {
	ptcl := flag.String("type", "tcp", "the protocol to use ")
	host := flag.String("host", "127.0.0.1", "The target host to scan")
	ports := flag.String("ports", "20-1000", "list of ports to scan ex : 80,443 or 0-100")
	timeout := flag.Int("timeout", 5, "Timeout in seconds for the connection")
	flag.Parse()

	portRange, err := parsePorts(*ports)

	var finalHost string
	if checkIP(*host) {
		finalHost = *host
	} else {
		if matched, _ := regexp.MatchString(`^([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`, *host); matched {
			ip, err := reverseLookup(*host)
			if err != nil {
				fmt.Printf("Erreur de reverse lookup pour l'URL %s : %v\n", *host, err)
				return
			}
			finalHost = ip
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

	scanAndGrab(finalHost, portRange, *timeout)
}