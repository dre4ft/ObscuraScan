package main

import (
    //"fmt"
    "net"
    "strings"
    "time"
)

// Structure pour les scripts de banner grabbing spécifiques
type ServiceGrabber struct {
    Port     int
    Protocol string
    Grabber  func(net.Conn) string
}

// Liste des scripts de banner grabbing pour ports bien connus
var serviceGrabbers = []ServiceGrabber{
    {
        Port:     21,   // FTP
        Protocol: "ftp",
        Grabber: func(conn net.Conn) string {
            buf := make([]byte, 1024)
            n, err := conn.Read(buf)
            if err == nil && n > 0 {
                banner := string(buf[:n])
                if strings.Contains(strings.ToLower(banner), "ftp") {
                    return strings.TrimSpace(banner)
                }
            }
            return "FTP server (generic banner)"
        },
    },
    {
        Port:     22,   // SSH
        Protocol: "ssh",
        Grabber: func(conn net.Conn) string {
            buf := make([]byte, 1024)
            n, err := conn.Read(buf)
            if err == nil && n > 0 {
                banner := string(buf[:n])
                if strings.Contains(strings.ToLower(banner), "ssh") {
                    return extractSSHVersion(banner)
                }
            }
            return "SSH server (generic banner)"
        },
    },
    {
        Port:     23,   // Telnet
        Protocol: "telnet",
        Grabber: func(conn net.Conn) string {
            buf := make([]byte, 1024)
            n, err := conn.Read(buf)
            if err == nil && n > 0 {
                banner := string(buf[:n])
                return strings.TrimSpace(banner)
            }
            return "Telnet server (no specific banner)"
        },
    },
    {
        Port:     25,   // SMTP
        Protocol: "smtp",
        Grabber: func(conn net.Conn) string {
            buf := make([]byte, 1024)
            n, err := conn.Read(buf)
            if err == nil && n > 0 {
                banner := string(buf[:n])
                if strings.Contains(strings.ToLower(banner), "smtp") {
                    return strings.TrimSpace(banner)
                }
            }
            return "SMTP server (generic banner)"
        },
    },
    {
        Port:     80,   // HTTP
        Protocol: "http",
        Grabber: func(conn net.Conn) string {
            req := "GET / HTTP/1.1\r\nHost: example.com\r\nConnection: close\r\n\r\n"
            _, err := conn.Write([]byte(req))
            if err != nil {
                return "HTTP server (connection error)"
            }

            buf := make([]byte, 1024)
            n, err := conn.Read(buf)
            if err == nil && n > 0 {
                banner := string(buf[:n])
                return extractHTTPServerHeader(banner)
            }
            return "HTTP server (no banner)"
        },
    },
    {
        Port:     443,  // HTTPS
        Protocol: "https",
        Grabber: func(conn net.Conn) string {
            req := "GET / HTTP/1.1\r\nHost: example.com\r\nConnection: close\r\n\r\n"
            _, err := conn.Write([]byte(req))
            if err != nil {
                return "HTTPS server (connection error)"
            }

            buf := make([]byte, 1024)
            n, err := conn.Read(buf)
            if err == nil && n > 0 {
                banner := string(buf[:n])
                return extractHTTPServerHeader(banner)
            }
            return "HTTPS server (no banner)"
        },
    },
    {
        Port:     443,  // MySQL
        Protocol: "mysql",
        Grabber: func(conn net.Conn) string {
            buf := make([]byte, 1024)
            n, err := conn.Read(buf)
            if err == nil && n > 0 {
                banner := string(buf[:n])
                if strings.Contains(strings.ToLower(banner), "mysql") {
                    return strings.TrimSpace(banner)
                }
            }
            return "MySQL server (generic banner)"
        },
    },
    {
        Port:     3306, // MySQL
        Protocol: "mysql",
        Grabber: func(conn net.Conn) string {
            buf := make([]byte, 1024)
            n, err := conn.Read(buf)
            if err == nil && n > 0 {
                banner := string(buf[:n])
                if strings.Contains(strings.ToLower(banner), "mysql") {
                    return strings.TrimSpace(banner)
                }
            }
            return "MySQL server (generic banner)"
        },
    },
}

// Fonction utilitaire pour extraire la version SSH
func extractSSHVersion(banner string) string {
    lines := strings.Split(banner, "\n")
    for _, line := range lines {
        if strings.Contains(line, "SSH") {
            return strings.TrimSpace(line)
        }
    }
    return "SSH server (version not found)"
}

// Fonction utilitaire pour extraire l'en-tête du serveur HTTP
func extractHTTPServerHeader(response string) string {
    lines := strings.Split(response, "\n")
    for _, line := range lines {
        if strings.HasPrefix(strings.ToLower(line), "server:") {
            return strings.TrimSpace(line)
        }
    }
    return "HTTP server (no server header)"
}

// Fonction pour obtenir le grabber spécifique à un port
func getServiceGrabber(port int) *ServiceGrabber {
    for _, sg := range serviceGrabbers {
        if sg.Port == port {
            return &sg
        }
    }
    return nil
}

// Fonction de banner grabbing améliorée
func advancedBannerGrab(conn net.Conn, port int) string {
    // Tentative avec un grabber spécifique
    grabber := getServiceGrabber(port)
    if grabber != nil {
        banner := grabber.Grabber(conn)
        if banner != "" {
            return banner
        }
    }
    
    // Fallback sur une méthode générique si aucun grabber spécifique n'est trouvé
    return genericBannerGrab(conn)
}

// Méthode générique de banner grabbing
func genericBannerGrab(conn net.Conn) string {
    conn.SetReadDeadline(time.Now().Add(5 * time.Second))
    
    buf := make([]byte, 4096)
    n, err := conn.Read(buf)
    if err == nil && n > 0 {
        banner := string(buf[:n])
        
        // Nettoyage basique
        banner = strings.TrimSpace(banner)
        
        if banner != "" {
            return banner
        }
    }
    
    return "No banner detected"
}

// Vous pouvez ajouter plus de scripts et de logiques ici