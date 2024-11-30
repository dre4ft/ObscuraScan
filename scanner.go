package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)


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
