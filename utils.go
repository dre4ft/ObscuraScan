package main

// 2. Importation des packages nécessaires
import (
	"net"
)


func checkIP(providedIp string)bool{
	parsedIP := net.ParseIP(providedIp)
	return parsedIP != nil
}

func checkProtocole(providedProtocol string)bool{
	return providedProtocol == "tcp" || providedProtocol == "udp"
}

func reverselookup(url string) []string{
	ips, err := net.LookupHost(url)
	if err != nil {
		fmt.Println("Erreur:", err)
		return
	}

	for _, ip := range ips {
		fmt.Println("IP:", ip)
	}
	return ips
}