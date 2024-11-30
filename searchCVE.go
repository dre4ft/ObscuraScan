package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

// Structure to store CVE data matching the NVD API response
type CVE struct {
	CVEID       string `json:"id"`
	Description string `json:"descriptions.0.value"`
	ExploitabilityScore string `json:"metrics.cvssMetricV3.0.cvssData.baseScore"`
}

// Fonction pour extraire le nom du service et la version de la bannière
func ExtractProtocolAndVersion(banner string) (string, string) {
	// Ignorer les bannières contenant "(no response)"
	if strings.Contains(banner, "(no response)") {
		return "", ""
	}

	// Liste des expressions régulières
	re1 := regexp.MustCompile(`(?P<Protocol>[a-zA-Z\-]+)[/\s](?P<Version>[0-9\.]+(?:[a-zA-Z\-0-9]+)?)`)
	re2 := regexp.MustCompile(`(?P<Protocol>[a-zA-Z\-]+)[_/](?P<Version>[0-9\.]+(?:[a-zA-Z\-0-9]+)?)`)
	re3 := regexp.MustCompile(`(?P<Protocol>[a-zA-Z\s\-]+)`)

	// Première expression régulière (Nom/Version ou Nom Version)
	match1 := re1.FindStringSubmatch(banner)
	if len(match1) > 0 {
		return match1[1], match1[2]
	}

	// Deuxième expression régulière (OpenSSH, SMB, etc.)
	match2 := re2.FindStringSubmatch(banner)
	if len(match2) > 0 {
		return match2[1], match2[2]
	}

	// Troisième expression régulière (cas sans version, par exemple SMB service)
	match3 := re3.FindStringSubmatch(banner)
	if len(match3) > 0 {
		return match3[1], "" // Pas de version
	}

	// Si aucune correspondance, retourner une chaîne vide
	return "", ""
}

// Recherche des CVEs en fonction du protocole et de la version
func SearchCVE(application, version string) []CVE {
	// Construire l'URL de recherche pour l'API NVD avec l'application et la version
	url := fmt.Sprintf("https://services.nvd.nist.gov/rest/json/cves/2.0?keywordSearch=%s+%s&resultsPerPage=5", application, version)

	// Faire une requête HTTP GET vers l'API NVD
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching CVEs:", err)
		return nil
	}
	defer resp.Body.Close()

	// Lire la réponse JSON
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}

	// Décoder la réponse JSON dans une structure Go
	var result struct {
		Vulnerabilities []struct {
			CVE CVE `json:"cve"`
		} `json:"vulnerabilities"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil
	}

	// Extraire les CVEs
	cves := []CVE{}
	for _, vuln := range result.Vulnerabilities {
		cves = append(cves, vuln.CVE)
	}

	// Retourner la liste des CVEs trouvés
	return cves
}

// Vérifier si une CVE est exploitable
func IsExploitable(cve CVE) string {
	// Si la CVE a un exploitabilityScore > 0, elle est considérée comme exploitable
	if cve.ExploitabilityScore != "" {
		return "yes"
	}
	return "no"
}

// Formater les résultats des CVEs pour un protocole/application donné
func FormatCVEResults(application, version string, cves []CVE) string {
	// Commencer à formater la chaîne avec le nom du protocole et de la version
	result := fmt.Sprintf("%s/%s :\n", application, version)
	if len(cves) > 0 {
		// Ajouter chaque CVE à la chaîne avec son état d'exploitabilité
		for _, cve := range cves {
			exploitability := IsExploitable(cve)
			result += fmt.Sprintf("- CVE %s exploitable: %s\n", cve.CVEID, exploitability)
		}
	}else {
		result += ("no CVE founded")
	}
	

	return result
}

// Wrapper : prendre une bannière, extraire le protocole/version, rechercher les CVEs et retourner les résultats formatés
func SearchAndFormat(banner string) string {
	// Extraire le protocole et la version de la bannière
	protocol, version := ExtractProtocolAndVersion(banner)

	// Si aucune information de protocole/version n'est extraite, retourner une erreur
	if protocol == "" || version == "" {
		return "Unable to extract protocol and version from banner."
	}

	// Rechercher les CVEs pour ce protocole/version
	cves := SearchCVE(protocol, version)

	// Retourner les résultats formatés
	return FormatCVEResults(protocol, version, cves)
}