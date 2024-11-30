package main

import (
	"fmt"
	"net"
	"time"
	 "sync"
	//"strings"     
   // "regexp"
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
    // Canal pour stocker les résultats
    results := make(map[string]string)
    
    // Canal pour synchroniser les résultats
    resultsChan := make(chan ScanResult, len(ports))
    
    // WaitGroup pour attendre la fin de tous les scans
    var wg sync.WaitGroup
    
    // Limiter le nombre de goroutines simultanées
    semaphore := make(chan struct{}, 100) // Limite à 100 connexions simultanées
    
    // Scanner chaque port en parallèle
    for _, port := range ports {
        wg.Add(1)
        semaphore <- struct{}{} // Acquérir un slot
        
        go func(port int) {
            defer wg.Done()
            defer func() { <-semaphore }() // Libérer le slot
            
            address := fmt.Sprintf("%s:%d", ip, port)
            conn, err := net.DialTimeout("tcp", address, time.Duration(timeout)*time.Second)
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
            resultsChan <- result
        }(port)
    }
    
    // Goroutine pour fermer le canal une fois tous les scans terminés
    go func() {
        wg.Wait()
        close(resultsChan)
    }()
    
    // Collecter les résultats
    for result := range resultsChan {
        protPort := fmt.Sprintf("%d/tcp", result.Port)
        if result.Banner != "" {
            results[protPort] = fmt.Sprintf("open\n%s", result.Banner)
        } else {
            results[protPort] = "open"
        }
    }
    
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
