package main

import (
	//"context"
    "fmt"
    "net"
    "sync"
    "time"
)



func grabBanner(conn net.Conn, port int) string {
    // Utilisez la nouvelle fonction advancedBannerGrab
    return advancedBannerGrab(conn, port)
}



// Structure pour stocker les résultats du scan
type ScanResult struct {
	Port   int
	State  string
	Banner string
}

// Fonction pour scanner une IP et ses ports en parallèle
func scan(ip string, ports []int, timeout int, grab bool) map[string]string {
	toReturn := make(map[string]string)
	var wg sync.WaitGroup // Crée un WaitGroup pour synchroniser les goroutines

	// Mutex pour protéger l'accès à la map partagée
	var mu sync.Mutex

	// Pour chaque port, lancer une goroutine
	for _, port := range ports {
		wg.Add(1) // Incrémente le compteur de goroutines
		go func(port int) {
			defer wg.Done() // Décrémente le compteur lorsque la goroutine termine

			address := fmt.Sprintf("%s:%d", ip, port)
			conn, err := net.DialTimeout("tcp", address, time.Duration(timeout)*time.Second)
			if err != nil {
				// Si le port est fermé ou inaccessible, on continue
				return
			}
			defer conn.Close()

			protPort := fmt.Sprintf("%d/tcp", port)

			// Récupère la bannière si demandé
			var result string
			if grab {
				banner := grabBanner(conn, port)
				result = fmt.Sprintf("open\n%s", banner)
			} else {
				result = "open"
			}

			// Verrouille l'accès à la map partagée pour l'écriture
			mu.Lock()
			toReturn[protPort] = result
			mu.Unlock()
		}(port)
	}

	// Attendez que toutes les goroutines aient terminé
	wg.Wait()

	return toReturn
}


/*

func scan(ip string, ports []int, timeout int, grab bool) map[string]string {
    toReturn := make(map[string]string)
    
    for _, port := range ports {
        address := fmt.Sprintf("%s:%d", ip, port)
        conn, err := net.DialTimeout("tcp", address, time.Duration(timeout)*time.Second)
        if err != nil {
            continue // Port fermé ou non accessible
        }
        defer conn.Close()

        protPort := fmt.Sprintf("%d/tcp", port)
        
        if grab {
            banner := grabBanner(conn, port)
            toReturn[protPort] = fmt.Sprintf("open\n%s", banner)
        } else {
            toReturn[protPort] = "open"
        }
    }
    
    return toReturn
}*/

// Convertit une map en une chaîne lisible.
func toString(input map[string]string) {
	for key, value := range input {
		cve := SearchAndFormat(value)
		fmt.Printf("%s : %s \n", key,value)
		fmt.Printf("Vulnerability scan : \n %s \n",cve)
	}
}

func scanIPRange(ips []string , ports []int , timeout int,grab bool){
	for _,ip := range ips {
		fmt.Printf("%s : \n", ip)
		toString(scan(ip,ports,timeout,grab))
	}
}
