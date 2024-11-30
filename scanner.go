package main

import (
	"context"
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
    Port    int
    State   string
    Banner  string
}

func parallelScan(ip string, ports []int, timeout int, grab bool) map[string]string {
    // Utilisation de canaux avec buffer
    results := make(map[string]string)
    resultsChan := make(chan ScanResult, len(ports))
    
    // Canal d'erreur et de synchronisation
    errChan := make(chan error, 1)
    
    // WaitGroup pour gérer les goroutines
    var wg sync.WaitGroup
    
    // Mutex pour sécuriser l'accès aux résultats
    var mu sync.Mutex
    
    // Limiter le nombre de goroutines simultanées
    semaphore := make(chan struct{}, 100)

    // Contexte avec timeout global
    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
    defer cancel()

    // Scanner chaque port en parallèle
    for _, port := range ports {
        wg.Add(1)
        semaphore <- struct{}{} // Acquérir un slot
        
        go func(port int) {
            defer wg.Done()
            defer func() { <-semaphore }() // Libérer le slot

            // Utilisation du contexte pour le timeout
            conn, err := (&net.Dialer{
                Timeout: time.Duration(timeout) * time.Second,
            }).DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", ip, port))
            
            if err != nil {
                return // Port fermé ou non accessible
            }
            defer conn.Close()

            result := ScanResult{
                Port:   port,
                State:  "open",
                Banner: "",
            }
            
            // Si grab est activé, tenter de récupérer la bannière
            if grab {
                result.Banner = grabBanner(conn, port)
            }
            
            // Envoyer le résultat au canal
            select {
            case resultsChan <- result:
            case <-ctx.Done():
                return
            }
        }(port)
    }
    
    // Goroutine pour fermer les canaux
    go func() {
        wg.Wait()
        close(resultsChan)
        close(errChan)
    }()

    // Collecter les résultats
    for result := range resultsChan {
        mu.Lock()
        protPort := fmt.Sprintf("%d/tcp", result.Port)
        if result.Banner != "" {
            results[protPort] = fmt.Sprintf("open\n%s", result.Banner)
        } else {
            results[protPort] = "open"
        }
        mu.Unlock()
    }
    toString(results)
    return results
}

// Remplacez votre fonction scan existante par cette nouvelle version
func scan(ip string, ports []int, timeout int, grab bool) map[string]string {
    return parallelScan(ip, ports, timeout, grab)
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
