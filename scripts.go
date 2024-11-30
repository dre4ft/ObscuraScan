package main

import (
    "fmt"
    "net"
    "strings"
    "time"
	"binary"
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
	{
        Port:     88,   // Kerberos
        Protocol: "kerberos",
        Grabber: func(conn net.Conn) string {
            // Tentative de récupération d'informations Kerberos
            buf := make([]byte, 1024)
            n, err := conn.Read(buf)
            if err == nil && n > 0 {
                // Analyse basique du protocole Kerberos
                if len(buf) >= 4 {
                    // Vérification des premiers octets typiques de Kerberos
                    length := binary.BigEndian.Uint32(buf[:4])
                    if length > 0 && length < 1024 {
                        return fmt.Sprintf("Kerberos server (packet length: %d)", length)
                    }
                }
                return "Kerberos server (generic response)"
            }
            return "Kerberos server (no response)"
        },
    },
    {
        Port:     135,  // Microsoft RPC
        Protocol: "msrpc",
        Grabber: func(conn net.Conn) string {
            // Tentative de récupération d'informations RPC
            buf := make([]byte, 1024)
            n, err := conn.Read(buf)
            if err == nil && n > 0 {
                // Analyse des données RPC
                if len(buf) >= 5 && buf[0] == 0x05 && buf[1] == 0x00 {
                    return "Microsoft RPC Endpoint Mapper"
                }
                return "RPC service (unidentified)"
            }
            return "RPC service (no response)"
        },
    },
    {
        Port:     139,  // NetBIOS
        Protocol: "netbios",
        Grabber: func(conn net.Conn) string {
            // Tentative de récupération d'informations NetBIOS
            buf := make([]byte, 1024)
            n, err := conn.Read(buf)
            if err == nil && n > 0 {
                // Analyse des données NetBIOS
                if len(buf) > 4 && buf[0] == 0xff && buf[1] == 0x53 {
                    return "NetBIOS Name Service"
                }
                return "NetBIOS service (unidentified)"
            }
            return "NetBIOS service (no response)"
        },
    },
    {
        Port:     389,  // LDAP
        Protocol: "ldap",
        Grabber: func(conn net.Conn) string {
            // Tentative de récupération d'informations LDAP
            buf := make([]byte, 1024)
            n, err := conn.Read(buf)
            if err == nil && n > 0 {
                banner := string(buf[:n])
                if strings.Contains(strings.ToLower(banner), "ldap") {
                    return strings.TrimSpace(banner)
                }
                return "LDAP server (generic response)"
            }
            return "LDAP server (no response)"
        },
    },
    {
        Port:     445,  // SMB
        Protocol: "smb",
        Grabber: func(conn net.Conn) string {
            // Tentative de récupération d'informations SMB
            // Envoi d'un paquet SMB minimal
            smbNegociate := []byte{
                0x00, 0x00, 0x00, 0x54, // Longueur du paquet
                0xff, 0x53, 0x4d, 0x42, // Signature SMB
                0x72,                   // Commande SMB Negotiate
                0x00, 0x00, 0x00, 0x00, // Status
                0x18, 0x01,             // Flags
                0x00, 0x00,             // Flags2
            }
            
            _, err := conn.Write(smbNegociate)
            if err != nil {
                return "SMB service (send error)"
            }
            
            buf := make([]byte, 1024)
            n, err := conn.Read(buf)
            if err == nil && n > 0 {
                // Vérification de la signature SMB
                if len(buf) >= 4 && buf[0] == 0x00 && buf[1] == 0x00 && buf[2] == 0x00 && buf[3] == 0x54 {
                    return "Microsoft SMB service"
                }
                return "SMB-like service (unidentified)"
            }
            return "SMB service (no response)"
        },
    },
    {
        Port:     1433, // Microsoft SQL Server
        Protocol: "mssql",
        Grabber: func(conn net.Conn) string {
            // Tentative de préparation d'une préconnexion SQL Server
            preLoginPacket := []byte{
                0x12, 0x01, 0x00, 0x34, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x15, 0x00, 0x06, 0x01, 
                0x00, 0x1b, 0x00, 0x01, 0x02, 0x00, 0x1c, 
                0x00, 0x0c, 0x03, 0x00, 0x1d, 0x00, 0x00, 
                0x00, 0x00, 0x04, 0x02, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00,
            }
            
            _, err := conn.Write(preLoginPacket)
            if err != nil {
                return "MSSQL service (send error)"
            }
            
            buf := make([]byte, 1024)
            n, err := conn.Read(buf)
            if err == nil && n > 0 {
                // Vérification de la réponse SQL Server
                if len(buf) >= 5 && buf[0] == 0x04 && buf[1] == 0x01 {
                    return "Microsoft SQL Server"
                }
                return "SQL-like service (unidentified)"
            }
            return "MSSQL service (no response)"
        },
    },
    {
        Port:     5432, // PostgreSQL
        Protocol: "postgresql",
        Grabber: func(conn net.Conn) string {
            // Paquet de négociation PostgreSQL
            postgresPacket := []byte{
                0x00, 0x00, 0x00, 0x4f, // Longueur
                0x00, 0x03, 0x00, 0x00, // Version du protocole
                0x75, 0x73, 0x65, 0x72, 0x00, // "user"
                0x70, 0x6f, 0x73, 0x74, 0x67, 0x72, 0x65, 0x73, 0x00, // "postgres"
                0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x00, // "database"
                0x74, 0x65, 0x73, 0x74, 0x00, // "test"
                0x00, // Fin
            }
            
            _, err := conn.Write(postgresPacket)
            if err != nil {
                return "PostgreSQL service (send error)"
            }
            
            buf := make([]byte, 1024)
            n, err := conn.Read(buf)
            if err == nil && n > 0 {
                // Vérification de la réponse PostgreSQL
                if len(buf) >= 1 && (buf[0] == 'R' || buf[0] == 'E') {
                    return "PostgreSQL database server"
                }
                return "PostgreSQL-like service (unidentified)"
            }
            return "PostgreSQL service (no response)"
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

