package main

// 2. Importation des packages nécessaires
import (
	"net"
	"fmt"
)


func checkIP(providedIp string)bool{
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil
}

func checkProtocole(providedProtocol string)bool{
	return providedProtocol == "tcp" || providedProtocol == "udp"
}