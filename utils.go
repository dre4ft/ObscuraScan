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