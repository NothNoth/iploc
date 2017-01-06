package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var (
	api = "http://ip-api.com/json/"
)

/*
{"as":"AS16509 Amazon.com, Inc.",
  "city":"Seattle",
  "country":"United States",
  "countryCode":"US",
  "isp":"Amazon Technologies",
  "lat":47.6103,
  "lon":-122.3341,
  "org":"Amazon.com",
  "query":"54.192.25.8",
  "region":"WA",
  "regionName":"Washington",
  "status":"success",
  "timezone":"America/Los_Angeles",
  "zip":"98101"}
*/

type Descriptor struct {
	As          string  `json:"as"`
	City        string  `json:"city"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Isp         string  `json:"isp"`
	Lat         float32 `json:"lat"`
	Lon         float32 `json:"lon"`
	Org         string  `json:"org"`
	Query       string  `json:"query"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`

	Status   string `json:"status"`
	Timezone string `json:"timezone"`

	Zip     string `json:"zip"`
	Reverse string `json:"reverse,omitempty"`
	Mobile  bool   `json:"mobile,omitempty"`
	Proxy   bool   `json:"proxy,omitempty"`
}

type queryError struct {
	Message string `json:"message"`
	Query   string `json:"query"`
	Status  string `json:"status"`
}

func processIP(ip string) (identification string) {

	resp, err := http.Get(fmt.Sprintf("%s%s", api, ip))

	if err != nil {
		return fmt.Sprintf("%s (query failed)", ip)
	}
	var desc Descriptor
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &desc)
	if err != nil {
		var errorDesc queryError
		err = json.Unmarshal(body, &errorDesc)
		if err != nil {
			return fmt.Sprintf("[%s:json failed]", ip)
		}

		return fmt.Sprintf("[%s:%s]", ip, errorDesc.Message)
	}

	if desc.Status == "fail" {
		return fmt.Sprintf("[%s:?/?/?]", ip)
	}

	if len(desc.Isp) == 0 {
		desc.Isp = "?"
	}
	if len(desc.CountryCode) == 0 {
		desc.CountryCode = "?"
	}
	if len(desc.City) == 0 {
		desc.City = "?"
	}

	return fmt.Sprintf("[%s:%s/%s/%s]", ip, desc.Isp, desc.CountryCode, desc.City)
}

func filter(line string) (output string) {
	re, err := regexp.Compile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
	if err != nil {
		panic("bad regexp")
	}

	//Extract all IPs from string
	allIPs := re.FindAllStringIndex(line, -1)

	if len(allIPs) == 0 {
		return line
	}

	//Patch (revers order to preserve indexes while patching)
	for z := len(allIPs) - 1; z >= 0; z-- {
		ipLoc := allIPs[z]
		ip := line[ipLoc[0]:ipLoc[1]]
		ipIdent := processIP(ip)
		line = strings.Replace(line, ip, ipIdent, 1)
	}

	return line
}

func main() {

	if len(os.Args) != 1 {
		fmt.Println(os.Args[0], "will try to detect ip from input text and replace ip with ip identification")
		fmt.Println("")
		fmt.Println("Usage:", os.Args[0], "<text>")
		fmt.Println("")
		fmt.Println("Typical uses:")
		fmt.Println("\tcat iplist.txt |", os.Args[0])
		fmt.Println("\tls |", os.Args[0])
		fmt.Println("\tcat /var/log/message |", os.Args[0])
		return
	}
	reader := bufio.NewReader(os.Stdin)

	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Print(filter(text))
	}
}
