package main

import (
	"fmt"
	"net"
	"time"
	"strings"     
    "regexp"
)


/*func scan (ip string,ports []int, timeout int , grab bool) map[string]string {
	
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
			conn.SetReadDeadline(time.Now().Add(30 * time.Second))

			// Utiliser un buffer pour lire la réponse
			var banner string
			buf := make([]byte, 1024)
			for {
				n, err := conn.Read(buf)
				if n > 0 {
					banner += string(buf[:n])
				}
				if err != nil {
					if err == io.EOF {
						break
					}
					banner = ""
					break
				}
			}

			if banner != "" {
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
}*/



func grabBanner(conn net.Conn, timeout time.Duration) string {
    conn.SetReadDeadline(time.Now().Add(timeout))
    
    // Buffer plus grand pour capturer plus d'informations
    buf := make([]byte, 4096)
    
    // Première tentative : lire directement
    n, err := conn.Read(buf)
    if err == nil && n > 0 {
        banner := string(buf[:n])
        
        // Nettoyage et filtrage de la bannière
        banner = cleanBanner(banner)
        
        if banner != "" {
            return banner
        }
    }
    
    // Deuxième tentative : essayer d'envoyer un caractère et relire
    if _, err := conn.Write([]byte("\r\n")); err == nil {
        n, err = conn.Read(buf)
        if err == nil && n > 0 {
            banner := string(buf[:n])
            banner = cleanBanner(banner)
            
            if banner != "" {
                return banner
            }
        }
    }
    
    return "Pas de banner trouvé"
}

func cleanBanner(banner string) string {
    // Nettoyer la bannière des caractères non imprimables
    banner = strings.TrimSpace(banner)
    
    // Filtrer les caractères non ASCII imprimables
    re := regexp.MustCompile(`[^\x20-\x7E\n\r]`)
    banner = re.ReplaceAllString(banner, "")
    
    // Limiter la longueur
    if len(banner) > 512 {
        banner = banner[:512]
    }
    
    // Supprimer les lignes vides
    var cleanLines []string
    for _, line := range strings.Split(banner, "\n") {
        trimmed := strings.TrimSpace(line)
        if trimmed != "" {
            cleanLines = append(cleanLines, trimmed)
        }
    }
    
    return strings.Join(cleanLines, "\n")
}

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
            banner := grabBanner(conn, 5*time.Second)
            toReturn[protPort] = fmt.Sprintf("open\n%s", banner)
        } else {
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
