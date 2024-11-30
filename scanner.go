package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

/*
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

func bannerGrab(portUp map[int]bool, ip string, timeout time.Duration) map[int]string {
	portBanner := make(map[int]string)

	// Boucle sur chaque port ouvert
	for port := range portUp {
		address := fmt.Sprintf("%s:%d", ip, port)
		conn, err := net.DialTimeout("tcp", address, timeout)
		if err != nil {
			// Si la connexion échoue, on l'ignore et on passe au port suivant
			fmt.Printf("Erreur de connexion au port %d : %v\n", port, err)
			continue
		}

		// Assurez-vous de fermer la connexion à la fin de la fonction
		defer conn.Close()

		// Pour le port 80 (HTTP), envoyer une requête GET
		if port == 80 {
			fmt.Fprintf(conn, "GET / HTTP/1.1\r\nHost: %s\r\nConnection: close\r\n\r\n", ip)
		}

		// Définir un délai de lecture (5 secondes)
		conn.SetReadDeadline(time.Now().Add(timeout))

		// Utiliser un buffer pour lire la réponse
		reader := bufio.NewReader(conn)
		var banner string
		// Lire plusieurs lignes si nécessaire (jusqu'à un maximum ou une nouvelle ligne)
		for {
			line, err := reader.ReadString('\n')
			banner += line
			if err != nil {
				if err.Error() != "EOF" {
					banner = ""
				}
				break
			}
		}

		if banner != "" {
			portBanner[port] = banner
		} else {
			// Si aucun banner n'est récupéré, mentionner que le port est ouvert mais sans réponse
			portBanner[port] = "Pas de banner trouvé"
		}
	}

	// Retourner les bannières récupérées pour chaque port
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
}*/



func scan (ip string,ports []int, timeout int , grab bool) map[string]string {
	
	toReturn := make(map[string]string)

	for _, port := range ports {
		isUp := false
		address := fmt.Sprintf("%s:%d", ip, port)

		conn, err := net.DialTimeout("tcp", address, time.Duration(timeout) * time.Second)
		if err != nil {
			continue // Si la connexion échoue, passer au port suivant
		} else {
			isUp = true
		}

		// Assurez-vous de fermer la connexion après l'utilisation
		defer conn.Close()

		protPort := fmt.Sprintf("tcp, %d", port)
		if grab && isUp {
			// Pour le port 80 (HTTP), envoyer une requête GET
			if port == 80 {
				fmt.Fprintf(conn, "GET / HTTP/1.1\r\nHost: %s\r\nConnection: close\r\n\r\n", ip)
			}

			// Définir un délai de lecture
			conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))

			// Utiliser un buffer pour lire la réponse
			reader := bufio.NewReader(conn)
			var banner string
			var hasAnErr = false 
			for {
				line, err := reader.ReadString('\n')
				banner += line
				if err != nil {
					hasAnErr = true 
					break
				}
			}

			if !hasAnErr {
				toReturn[protPort] = fmt.Sprintf("open\n%s", banner)
			} else {
				toReturn[protPort] = "open\nPas de banner trouvé"
			}
		} else if isUp {
			toReturn[protPort] = "open"
		}
	}

	toString(toReturn)
	return toReturn
}

// Convertit une map en une chaîne lisible.
func toString(input map[string]string) {
	for key, value := range input {
		fmt.Printf("%s : %s \n", key,value )
	}
}

func scanIPRange(ips []string , ports []int , timeout int,grab bool){
	for _,ip := range ips {
		fmt.Printf("%s : \n", ip)
		scan(ip,ports,timeout,grab)
	}
}
