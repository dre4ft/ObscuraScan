package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

// Structure pour stocker les données de CVE
type CVE struct {
	CVEID       string `json:"CVE_data_meta.CVE_ID"`
	Description string `json:"description.description_data"`
	// Impact spécifique qui est un indicateur d'exploitabilité (Simplification)
	ExploitabilityScore string `json:"impact.baseMetricV2.exploitabilityScore"`
}

// Fonction pour extraire le nom du service et la version de la bannière
func ExtractProtocolAndVersion(banner string) (string, string) {
	// Ignorer les bannières contenant "(no response)"
	if strings.Contains(banner, "(no response)") {
		return "", ""
	}

	// Liste des expressions régulières
	// Chercher des formats comme "Nom/Version", "Nom Version"
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
	// Construire l'URL de recherche pour l'API NIST avec l'application et la version
	url := fmt.Sprintf("https://services.nist.gov/rest/v2/cve/?keyword=%s %s", application, version)

	// Faire une requête HTTP GET vers l'API NIST
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching CVEs:", err)
		return nil
	}
	defer resp.Body.Close()

	// Lire la réponse JSON
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}

	// Décoder la réponse JSON dans une structure Go
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil
	}

	// Extraire les CVEs
	cves := []CVE{}
	if data, exists := result["result"]; exists {
		if cveItems, exists := data.(map[string]interface{})["CVE_Items"]; exists {
			cveList := cveItems.([]interface{})
			for _, item := range cveList {
				cveData := item.(map[string]interface{})
				cveID := cveData["cve"].(map[string]interface{})["CVE_data_meta"].(map[string]interface{})["CVE_ID"].(string)
				description := cveData["cve"].(map[string]interface{})["description"].(map[string]interface{})["description_data"].([]interface{})[0].(map[string]interface{})["value"].(string)

				// On ne prend que les CVEs avec un exploitabilityScore
				var exploitabilityScore string
				if exploitability, exists := cveData["impact"].(map[string]interface{})["baseMetricV2"].(map[string]interface{})["exploitabilityScore"]; exists {
					exploitabilityScore = fmt.Sprintf("%v", exploitability)
				}

				cves = append(cves, CVE{
					CVEID:       cveID,
					Description: description,
					ExploitabilityScore: exploitabilityScore,
				})
			}
		}
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
	result := fmt.Sprintf("%s : %s :\n", application, version)

	// Ajouter chaque CVE à la chaîne avec son état d'exploitabilité
	for _, cve := range cves {
		exploitability := IsExploitable(cve)
		result += fmt.Sprintf("- CVE %s exploitable: %s\n", cve.CVEID, exploitability)
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