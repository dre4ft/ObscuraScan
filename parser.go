package main

// 2. Importation des packages nécessaires
import (
	"errors"
    "fmt"
    "strconv"
    "strings"
)
   
func parsePorts(input string) ([]int, error) {
    input = strings.TrimSpace(input) // Supprimer les espaces autour de l'entrée

    // Vérifier si c'est un range (ex: "0-100")
    if strings.Contains(input, "-") {
        parts := strings.Split(input, "-")
        if len(parts) != 2 {
            return nil, errors.New("invalid range format")
        }

        start, err1 := strconv.Atoi(parts[0])
        end, err2 := strconv.Atoi(parts[1])
        if err1 != nil || err2 != nil || start > end || end > 65535 {
            return nil, errors.New("invalid range values")
        }

        // Générer la liste des ports pour le range
        ports := make([]int, 0)
        for i := start; i <= end; i++ {
            ports = append(ports, i)
        }
        return ports, nil
    }

    // Vérifier si c'est une liste (ex: "23,24,80")
    if strings.Contains(input, ",") {
        parts := strings.Split(input, ",")
        ports := make([]int, 0, len(parts))
        for _, part := range parts {
            port, err := strconv.Atoi(strings.TrimSpace(part))
            if err != nil || port < 0 || port > 65535 {
                return nil, fmt.Errorf("invalid port: %s", part)
            }
            ports = append(ports, port)
        }
        return ports, nil
    }

    // Traiter comme un seul port (ex: "80")
    port, err := strconv.Atoi(input)
    if err != nil || port < 0 || port > 65535 {
        return nil, fmt.Errorf("invalid port: %s", input)
    }

    return []int{port}, nil
}

