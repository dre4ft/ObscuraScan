package main

import (
	"flag"
	"fmt"
	"regexp"
)

func main() {
	ptcl := flag.String("type", "tcp", "the protocol to use ")
	host := flag.String("ip", "127.0.0.1", "The target host to scan")
	ports := flag.String("p", "20-1000", "list of ports to scan ex : 80,443 or 0-100")
	timeout := flag.Int("t", 5, "Timeout in seconds for the connection")
	grab := flag.Bool("g",false,"grab the banner for each open port")
	flag.Parse()

	portRange, err := parsePorts(*ports)

	var singleHost string
	var multiHost []string 
	var is_singleHost bool = false 

	if checkIP(*host) {
		singleHost = *host
		is_singleHost = true
	} else {
		if matched, _ := regexp.MatchString(`^([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`, *host); matched {
			ip ,err:= reverselookup(*host)
			if err != nil {
				fmt.Printf("Error ip lookup  %s : %v\n", *host, err)
				return
			}
			multiHost = ip
		} else {
			fmt.Println("Error : incorrect URL type")
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
		toString(scan(singleHost, portRange, *timeout,*grab ))
	} else {
		scanIPRange(multiHost, portRange,*timeout,*grab )
	}

	
}