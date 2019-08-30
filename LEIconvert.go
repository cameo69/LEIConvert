package main

/*
LEIconvert by Dietrich

Converts GLEIF Golden Copy (LEIs) to csv file

Allows for filtering for specific LEIs

based on this article https://blog.singleton.io/posts/2012-06-19-parsing-huge-xml-files-with-go/
and this sample code https://github.com/dps/go-xml-parse/blob/master/go-xml-parse.go

GLEIF data available here: https://www.gleif.org/en/lei-data/gleif-golden-copy/download-the-golden-copy#/
look for LEI-CDF v2.1 xml files
*/

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"time"
)

var inputFile = flag.String("in", "gleif-goldencopy-lei2-golden-copy.xml", "input file path")
var outFile = flag.String("out", "found.csv", "output file path")
var filterFile = flag.String("filter", "filter.txt", "file with all LEI to filter for")
var missingFile = flag.String("missing", "missing.txt", "file that gets all missing LEIs")
var doDebug = flag.Bool("debug", false, "do a short run for debugging")

const lapVersion = "0.10"

// Here is an example from the XML
/*
<lei:LEIRecord xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
               xmlns:lei="http://www.gleif.org/data/schema/leidata/2016"
               xmlns:gleif="http://www.gleif.org/data/schema/golden-copy/extensions/1.0">
   <lei:LEI>001GPB6A9XPE8XJICC14</lei:LEI>
   <lei:Entity>
      <lei:LegalName>Fidelity Advisor Series I - Fidelity Advisor Leveraged Company Stock Fund</lei:LegalName>
      <lei:LegalAddress>
         <lei:FirstAddressLine>C/O Eqiuty Portfolio Growth</lei:FirstAddressLine>
         <lei:AdditionalAddressLine>82 Devonshire Street</lei:AdditionalAddressLine>
         <lei:City>Boston</lei:City>
         <lei:Region>US-MA</lei:Region>
         <lei:Country>US</lei:Country>
         <lei:PostalCode>02109</lei:PostalCode>
      </lei:LegalAddress>
      <lei:HeadquartersAddress>
         <lei:FirstAddressLine>245 Summer Street</lei:FirstAddressLine>
         <lei:City>Boston</lei:City>
         <lei:Region>US-MA</lei:Region>
         <lei:Country>US</lei:Country>
         <lei:PostalCode>02210</lei:PostalCode>
      </lei:HeadquartersAddress>
      <lei:RegistrationAuthority>
         <lei:RegistrationAuthorityID>RA888888</lei:RegistrationAuthorityID>
         <lei:RegistrationAuthorityEntityID>S000005113</lei:RegistrationAuthorityEntityID>
      </lei:RegistrationAuthority>
      <lei:LegalJurisdiction>US</lei:LegalJurisdiction>
      <lei:EntityCategory>FUND</lei:EntityCategory>
      <lei:LegalForm>
         <lei:EntityLegalFormCode>9999</lei:EntityLegalFormCode>
         <lei:OtherLegalForm>OTHER</lei:OtherLegalForm>
      </lei:LegalForm>
      <lei:EntityStatus>ACTIVE</lei:EntityStatus>
   </lei:Entity>
   <lei:Registration>
      <lei:InitialRegistrationDate>2012-11-29T16:33:00.000Z</lei:InitialRegistrationDate>
      <lei:LastUpdateDate>2018-07-10T21:31:00.000Z</lei:LastUpdateDate>
      <lei:RegistrationStatus>ISSUED</lei:RegistrationStatus>
      <lei:NextRenewalDate>2019-06-18T14:31:00.000Z</lei:NextRenewalDate>
      <lei:ManagingLOU>EVK05KS7XY1DEII3R011</lei:ManagingLOU>
      <lei:ValidationSources>FULLY_CORROBORATED</lei:ValidationSources>
   </lei:Registration>
   <lei:Extension>
      <gleif:Geocoding>
         <gleif:original_address>245 Summer Street, 02210, Boston, US-MA, US</gleif:original_address>
         <gleif:relevance>0.92</gleif:relevance>
         <gleif:match_type>pointAddress</gleif:match_type>
         <gleif:lat>42.3514</gleif:lat>
         <gleif:lng>-71.05385</gleif:lng>
         <gleif:geocoding_date>2017-10-23T19:14:11</gleif:geocoding_date>
         <gleif:bounding_box>TopLeft.Latitude: 42.3525242, TopLeft.Longitude: -71.0553711, BottomRight.Latitude: 42.3502758, BottomRight.Longitude: -71.0523289</gleif:bounding_box>
         <gleif:match_level>houseNumber</gleif:match_level>
         <gleif:formatted_address>245 Summer St, Boston, MA 02210, United States</gleif:formatted_address>
         <gleif:mapped_location_id>NT_PYMT6GOD3rrAC9q2Al5jZB_yQTN</gleif:mapped_location_id>
         <gleif:mapped_street>Summer St</gleif:mapped_street>
         <gleif:mapped_housenumber>245</gleif:mapped_housenumber>
         <gleif:mapped_postalcode>02210</gleif:mapped_postalcode>
         <gleif:mapped_city>Boston</gleif:mapped_city>
         <gleif:mapped_district>Downtown Boston</gleif:mapped_district>
         <gleif:mapped_state>MA</gleif:mapped_state>
         <gleif:mapped_country>USA</gleif:mapped_country>
      </gleif:Geocoding>
      <gleif:Geocoding>
         <gleif:original_address>82 Devonshire Street, 02109, Boston, US-MA, US</gleif:original_address>
         <gleif:relevance>0.92</gleif:relevance>
         <gleif:match_type>pointAddress</gleif:match_type>
         <gleif:lat>42.35786</gleif:lat>
         <gleif:lng>-71.05691</gleif:lng>
         <gleif:geocoding_date>2017-10-23T21:12:15</gleif:geocoding_date>
         <gleif:bounding_box>TopLeft.Latitude: 42.3589842, TopLeft.Longitude: -71.0584313, BottomRight.Latitude: 42.3567358, BottomRight.Longitude: -71.0553887</gleif:bounding_box>
         <gleif:match_level>houseNumber</gleif:match_level>
         <gleif:formatted_address>82 Devonshire St, Boston, MA 02109, United States</gleif:formatted_address>
         <gleif:mapped_location_id>NT_ljFPoKOQ8PLHLphJj2Tx1C_4ID</gleif:mapped_location_id>
         <gleif:mapped_street>Devonshire St</gleif:mapped_street>
         <gleif:mapped_housenumber>82</gleif:mapped_housenumber>
         <gleif:mapped_postalcode>02109</gleif:mapped_postalcode>
         <gleif:mapped_city>Boston</gleif:mapped_city>
         <gleif:mapped_district>Downtown Boston</gleif:mapped_district>
         <gleif:mapped_state>MA</gleif:mapped_state>
         <gleif:mapped_country>USA</gleif:mapped_country>
      </gleif:Geocoding>
   </lei:Extension>
</lei:LEIRecord>
*/

//tLEIRecord Represents top level node of the interesting part in the XML file
type tLEIRecord struct {
	LEIRecord    string         `xml:"LEIRecord"`
	LEI          string         `xml:"LEI"`
	EntityList   []tEntity      `xml:"Entity"`
	Registration *tRegistration `xml:"Registration"`
}

type tEntity struct {
	XMLName             xml.Name  `xml:"Entity"`
	LegalName           string    `xml:"LegalName"`
	LegalAddress        *tAddress `xml:"LegalAddress"`
	HeadquartersAddress *tAddress `xml:"HeadquartersAddress"`
}

type tAddress struct {
	XMLName               xml.Name
	FirstAddressLine      string `xml:"FirstAddressLine"`
	AdditionalAddressLine string `xml:"AdditionalAddressLine"`
	City                  string `xml:"City"`
	Region                string `xml:"Region"`
	Country               string `xml:"Country"`
	PostalCode            string `xml:"PostalCode"`
}

type tRegistration struct {
	XMLName                 xml.Name `xml:"Registration"`
	InitialRegistrationDate string   `xml:"InitialRegistrationDate"`
	LastUpdateDate          string   `xml:"LastUpdateDate"`
	RegistrationStatus      string   `xml:"RegistrationStatus"`
	NextRenewalDate         string   `xml:"NextRenewalDate"`
	ManagingLOU             string   `xml:"ManagingLOU"`
	ValidationSources       string   `xml:"ValidationSources"`
}

func writeHeader(writer *bufio.Writer) {
	writer.WriteString("LEI;LegalName;LegalAddress_FirstAddressLine;LegalAddress_AdditionalAddressLine;LegalAddress_City;LegalAddress_Region;LegalAddress_Country;LegalAddress_PostalCode")
	writer.WriteString(";HeadquartersAddress_FirstAddressLine;HeadquartersAddress_AdditionalAddressLine;HeadquartersAddress_City;HeadquartersAddress_Region;HeadquartersAddress_Country")
	writer.WriteString(";HeadquartersAddress_PostalCode")

	writer.WriteString(";InitialRegistrationDate;LastUpdateDate;RegistrationStatus;NextRenewalDate;ManagingLOU;ValidationSources\n")
}
func writeRegistration(writer *bufio.Writer, regis *tRegistration) {
	writer.WriteString("\"" + regis.InitialRegistrationDate)
	writer.WriteString("\";\"" + regis.LastUpdateDate)
	writer.WriteString("\";\"" + regis.RegistrationStatus)
	writer.WriteString("\";\"" + regis.NextRenewalDate)
	writer.WriteString("\";\"" + regis.ManagingLOU)
	writer.WriteString("\";\"" + regis.ValidationSources)
	writer.WriteString("\"")
}
func writeAddress(writer *bufio.Writer, address *tAddress) {
	writer.WriteString("\"" + address.FirstAddressLine)
	writer.WriteString("\";\"" + address.AdditionalAddressLine)
	writer.WriteString("\";\"" + address.City)
	writer.WriteString("\";\"" + address.Region)
	writer.WriteString("\";\"" + address.Country)
	writer.WriteString("\";\"" + address.PostalCode)
	writer.WriteString("\"")
}

func writeLine(writer *bufio.Writer, p *tLEIRecord) {
	writer.WriteString("\"" + p.LEI + "\";\"" + p.EntityList[0].LegalName + "\";")

	writeAddress(writer, p.EntityList[0].LegalAddress)
	writer.WriteString(";")
	writeAddress(writer, p.EntityList[0].HeadquartersAddress)
	writer.WriteString(";")
	writeRegistration(writer, p.Registration)
	writer.WriteString("\n")
}

type myMap map[string]int

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) (myMap, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	//create map
	m := make(myMap)

	//var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		m[scanner.Text()] = 1 //1 means exist
	}

	return m, scanner.Err()
}

func matchFilter(filterLEIs myMap, LEI string) int {
	return filterLEIs[LEI]
}

func printStatus(totalLines int, totalWritten int, filterCount int, bytesTotal int64, bytesWritten int64, started time.Time) {
	ratio := float64(bytesWritten) / float64(bytesTotal)
	past := float64(time.Since(started))
	totalDur := time.Duration(past / ratio)
	estimated := started.Add(totalDur)
	remaining := time.Duration(time.Until(estimated)).Round(time.Second)

	if filterCount == 0 {
		fmt.Printf("\rRemaining: %v   Processed and written records: %d (%.1f%%)               ", remaining, totalWritten, ratio*100)
	} else {
		var prc float32
		prc = float32(totalWritten) / float32(filterCount) * 100
		fmt.Printf("\rRemaining: %v   Processed records: %d (%.1f%%)   Found and written records: %d/%d (%.1f%%)               ", remaining, totalLines, ratio*100, totalWritten, filterCount, prc)
	}
}

func remainingLEICount(m *myMap) int {
	zaehl := 0
	for _, v := range *m {
		//fmt.Printf("%s -> %s\n", k, v)
		if v == 1 {
			zaehl++
		}
	}
	return zaehl
}

func remainingLEI(m *myMap) myMap {
	newM := make(myMap)
	for k, v := range *m {
		//fmt.Printf("%s -> %s\n", k, v)
		if v == 1 {
			newM[k] = 1
		}
	}
	return newM
}

func writeRemainingLEI(m *myMap) {
	if *missingFile != "" {
		misFile, err := os.Create(*missingFile)
		if err != nil {
			fmt.Println("Error opening missing file:", err)
			return
		}
		defer misFile.Close()

		//write LEIs
		for k := range *m {
			fmt.Fprintln(misFile, k)
		}
	}
}

/*
#
#
#
#
#
#
#
#
#
#
*/

func main() {
	started := time.Now()

	fmt.Printf("LEI Golden Copy XML converter v%s\n", lapVersion)
	flag.Parse()

	var filterLEIs myMap
	var err error

	//enable below 2 comments to use the test-filter-list
	//filterTemp := "LEIs_for_filter.txt"
	//filterFile = &filterTemp
	if len(*filterFile) > 0 {
		filterLEIs, err = readLines(*filterFile)
		if err != nil {
			//log.Fatalf("readLines: %s", err)
			fmt.Printf("filter: %s: %s\n", err, *inputFile)
			filterLEIs = make(myMap)
		}
		fmt.Println("Filter LEIs: ", len(filterLEIs))
	} else {
		fmt.Println("Running without filter")
	}

	xmlFile, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("in: %s: %s\n", err, *inputFile)
		return
	}
	defer xmlFile.Close()

	xmlStat, err := xmlFile.Stat()
	xmlSize := xmlStat.Size()
	//fmt.Printf("The file is %d bytes long\n", xmlSize)

	decoder := xml.NewDecoder(xmlFile)

	totalLines := 0
	totalWritten := 0
	var inElement string

	outFile, err := os.Create(*outFile)
	if err != nil {
		fmt.Println("Error opening output file:", err)
		return
	}
	writer := bufio.NewWriter(outFile)
	defer outFile.Close()

	//write header
	writeHeader(writer)

	getOut := false
	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		// Inspect the type of the token just read.
		switch se := t.(type) {
		case xml.StartElement:
			// If we just read a StartElement token
			inElement = se.Name.Local
			// ...and its name is "page"
			if inElement == "LEIRecord" {
				var p tLEIRecord
				// decode a whole chunk of following XML into the
				// variable p which is a Page (see above)
				decoder.DecodeElement(&p, &se)

				if len(filterLEIs) == 0 {
					writeLine(writer, &p)
					totalWritten++
				} else {
					if matchFilter(filterLEIs, p.LEI) > 0 {
						writeLine(writer, &p)
						filterLEIs[p.LEI] = 2 //written
						totalWritten++
					}
				}

				totalLines++

				if totalLines%10000 == 0 {
					writer.Flush()
					printStatus(totalLines, totalWritten, len(filterLEIs), xmlSize, decoder.InputOffset(), started)
				}

				//way to get out for debugging
				if *doDebug && totalLines >= 20000 {
					getOut = true
				}

			}
		default:
		}

		if getOut {
			break
		}

		if len(filterLEIs) != 0 && totalWritten == len(filterLEIs) {
			//fmt.Println("\n", remainingLEICount(&filterLEIs))
			if remainingLEICount(&filterLEIs) == 0 {
				break
			}
		}

	}

	writer.Flush()
	printStatus(totalLines, totalWritten, len(filterLEIs), xmlSize, decoder.InputOffset(), started)
	fmt.Println()

	if len(filterLEIs) > 0 {
		if remainingLEICount(&filterLEIs) == 0 {
			fmt.Println("All LEI from filter list found")
		} else {
			fmt.Println("Not all LEI could be found:")
			missing := remainingLEI(&filterLEIs)
			for k := range missing {
				fmt.Println(k)
			}
			writeRemainingLEI(&missing)
		}
	} else if *doDebug {
		fmt.Println("CAREFUL, DEBUG MODE, file may be incomplete")
	} else {
		fmt.Println("File completely processed")
	}

	ende := time.Now()
	fmt.Println("Total processing time: ", ende.Sub(started).Round(time.Second))
}
