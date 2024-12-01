package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type CVE struct {
	CVEID       string `json:"id"`
	Description string `json:"descriptions.0.value"`
	ExploitabilityScore string `json:"metrics.cvssMetricV3.0.cvssData.baseScore"`
}


func ExtractProtocolAndVersion(banner string) (string, string) {

	if strings.Contains(banner, "(no response)") || strings.Contains(banner, "(generic banner)") ||  strings.Contains(banner, "(no specific banner)") ||  strings.Contains(banner,"(unidentified)"){
		return "", ""
	} 


	re1 := regexp.MustCompile(`(?P<Protocol>[a-zA-Z\-]+)[/\s](?P<Version>[0-9\.]+(?:[a-zA-Z\-0-9]+)?)`)
	re2 := regexp.MustCompile(`(?P<Protocol>[a-zA-Z\-]+)[_/](?P<Version>[0-9\.]+(?:[a-zA-Z\-0-9]+)?)`)
	re3 := regexp.MustCompile(`(?P<Protocol>[a-zA-Z\s\-]+)`)


	match1 := re1.FindStringSubmatch(banner)
	if len(match1) > 0 {
		return match1[1], match1[2]
	}


	match2 := re2.FindStringSubmatch(banner)
	if len(match2) > 0 {
		return match2[1], match2[2]
	}


	match3 := re3.FindStringSubmatch(banner)
	if len(match3) > 0 {
		return match3[1], "" 
	}

	return "", ""
}


func SearchCVE(application, version string) []CVE {

	url := fmt.Sprintf("https://services.nvd.nist.gov/rest/json/cves/2.0?keywordSearch=%s+%s&resultsPerPage=5", application, version)


	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching CVEs:", err)
		return nil
	}
	defer resp.Body.Close()


	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}


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


	cves := []CVE{}
	for _, vuln := range result.Vulnerabilities {
		cves = append(cves, vuln.CVE)
	}


	return cves
}


func IsExploitable(cve CVE) string {

	if cve.ExploitabilityScore != "" {
		return "yes"
	}
	return "no"
}


func FormatCVEResults(application, version string, cves []CVE) string {

	result := fmt.Sprintf("%s/%s :\n", application, version)
	if len(cves) > 0 {
		for _, cve := range cves {
			exploitability := IsExploitable(cve)
			result += fmt.Sprintf("- CVE %s exploitable: %s\n", cve.CVEID, exploitability)
		}
	}else {
		result += ("no CVE found")
	}
	

	return result
}


func SearchAndFormat(banner string) string {

	protocol, version := ExtractProtocolAndVersion(banner)


	if protocol == "" || version == "" {
		return "Unable to extract protocol and version from banner."
	}

	cves := SearchCVE(protocol, version)


	return FormatCVEResults(protocol, version, cves)
}