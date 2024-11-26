package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)


// Vérifie quels ports sont ouverts.
func whosup(ip string, portRange []int, timeout int) map[int]bool {
	toReturn := make(map[int]bool)
	duration := time.Duration(timeout) * time.Second
	for _, port := range portRange {
		toReturn[port] = scan("tcp", ip, port, duration)
	}
	return toReturn
}

// Vérifie si un port spécifique est ouvert.
func scan(protocol string, ip string, port int, timeout time.Duration) bool {
	ipPort := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout(protocol, ipPort, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// Supprime les entrées `false` d'une map.
func removeDown(input map[int]bool) map[int]bool {
	result := make(map[int]bool)
	for key, value := range input {
		if value {
			result[key] = value
		}
	}
	return result
}

// Effectue le banner grabbing sur les ports ouverts.
func bannerGrab(portUp map[int]bool, ip string, timeout time.Duration) map[int]string {
	portBanner := make(map[int]string)

	for port := range portUp {
		address := fmt.Sprintf("%s:%d", ip, port)
		conn, err := net.DialTimeout("tcp", address, timeout)
		if err != nil {
			fmt.Printf("Erreur de connexion au port %d : %v\n", port, err)
			continue
		}

		defer conn.Close()

		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		reader := bufio.NewReader(conn)
		banner, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Erreur de lecture sur le port %d : %v\n", port, err)
			continue
		}

		portBanner[port] = banner
	}
	return portBanner
}

// Combine les fonctions pour scanner et récupérer les bannières.
func scanAndGrab(ip string, portRange []int, timeout int) {
	openPorts := removeDown(whosup(ip, portRange, timeout))
	if len(openPorts)>0{
		banners := bannerGrab(openPorts, ip, time.Duration(timeout)*time.Second)
		for port, banner := range banners {
			fmt.Printf("tcp:%d : open\n%s\n\n", port, banner)
		}
	} else{
		fmt.Printf("no open port")
	}
}

// Convertit une map en une chaîne lisible.
func toString(input map[int]bool) {
	for key, value := range input {
		if value {
			fmt.Printf("%d : open\n", key)
		} else {
			fmt.Printf("%d : closed\n", key)
		}
	}
}
