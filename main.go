package main

import (
	"encoding/csv"
	"github.com/gocolly/colly"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	fname := "data.csv"
	file, err := os.Create(fname)
	if err != nil {
		log.Fatalf("Could not create file %s. Error: %v", fname, err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	c := colly.NewCollector(colly.AllowedDomains("internshala.com"))
	//c := colly.NewCollector()

	c.OnHTML(".internship_meta", func (e * colly.HTMLElement) {
		companyElement := e.ChildText(".company")
		location := e.ChildText("#location_names")
		otherDetails := e.ChildText(".item_body")
		title, company := extractCompanyAndJobTitle(companyElement)
		startDate, duration, stipend, applyBy, partTimeAllowed := extractOtherDetails(otherDetails)
		writer.Write([]string {
			title,
			company,
			location,
			startDate,
			duration,
			stipend,
			applyBy,
			partTimeAllowed,
		})
	})
	c.Visit("https://internshala.com/internships/page-1")
	log.Print("Scraping Complete")
}

func formatString(str string) string {
	regex := regexp.MustCompile(`\s\s+\n*`)
	formattedString := string(regex.ReplaceAll([]byte(str), []byte("|")))
	return formattedString
}

func extractCompanyAndJobTitle(companyElement string) (string, string) {
	companyElement = formatString(companyElement)
	split := strings.Split(companyElement, "|")
	jobTitle := split[0]
	companyName := split[1]
	return jobTitle, companyName
}

func extractOtherDetails(otherDetails string) (string, string, string, string, string) {
	otherDetails = formatString(otherDetails)
	split := strings.Split(otherDetails, "|")
	startDate:= split[0]
	duration := split[1]
	stipend := split[2]
	applyBy := split[3]
	partTimeFlag := "false"

	if strings.Contains(startDate, "Starts") {
		startDate = "Immediately"
	}
	if strings.Contains(applyBy, "Part time allowed") {
		applyBy = strings.ReplaceAll(applyBy, "Part time allowed", "")
		partTimeFlag = "true"
	}
	return startDate, duration, stipend, applyBy, partTimeFlag
}