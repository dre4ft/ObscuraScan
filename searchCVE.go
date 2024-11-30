package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

// CVEEntry représente une entrée CVE renvoyée par l'API du NIST
type CVEEntry struct {
	CveID      string `json:"id"`
	IsExploitable bool `json:"isExploitable"`
}

// CVEResponse représente la structure JSON de la réponse de l'API
type CVEResponse struct {
	Vulnerabilities []struct {
		CVE CVEEntry `json:"cve"`
	} `json:"vulnerabilities"`
}

// ExtractProtocolAndVersion extrait le protocole/l'application et la version d'une bannière
func ExtractProtocolAndVersion(banner string) (string, string) {
	// Exemple : Apache/2.4.51
	re := regexp.MustCompile(`(?P<Protocol>[a-zA-Z0-9\-_]+)/(?P<Version>[0-9]+\.[0-9]+(?:\.[0-9]+)?)`)
	match := re.FindStringSubmatch(banner)

	if len(match) == 0 {
		fmt.Println("Error: Unable to parse banner.")
		return "", ""
	}

	// Récupérer les noms des groupes
	protocol := match[1]
	version := match[2]

	return protocol, version
}

// BuildCPE construit un CPE standardisé à partir du protocole et de la version
func BuildCPE(protocol, version string) string {
	return fmt.Sprintf("cpe:2.3:a:%s:%s:%s:*:*:*:*:*:*:*", strings.ToLower(protocol), strings.ToLower(protocol), version)
}

// SearchCVEs interroge l'API du NIST pour récupérer les CVEs d'un protocole/version spécifique
func SearchCVEs(protocol, version string) []CVEEntry {
	cpe := BuildCPE(protocol, version)
	url := fmt.Sprintf("https://services.nvd.nist.gov/rest/json/cves/2.0?cpeName=%s&isVulnerable=true", cpe)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching CVEs: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return nil
	}

	var response CVEResponse
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return nil
	}

	var cves []CVEEntry
	for _, vuln := range response.Vulnerabilities {
		cves = append(cves, vuln.CVE)
	}

	return cves
}

// CondenseResults formate les résultats des CVE en une chaîne lisible
func CondenseResults(protocol, version string, cves []CVEEntry) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s : %s :\n", protocol, version))
	for _, cve := range cves {
		sb.WriteString(fmt.Sprintf("- CVE %s exploitable: %t\n", cve.CveID, cve.IsExploitable))
	}

	return sb.String()
}
