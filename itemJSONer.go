package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/gocolly/colly/v2"
)

const URL = "https://kol.coldfront.net/newmarket/"

// Another website option to get HTML from.
const URL2 = "https://g1wjmf0i0h.execute-api.us-east-2.amazonaws.com/default/itemindex"
const htmlFile = "newmarket.html"
const htmlFilePath = "file://./" + htmlFile
const jsonFile = "items.json"

var updateFileFlag = flag.Bool("u", false, "bool to get/update HTML file from URL")

type item struct {
	Name string `json:"name"`
	ID   string `json:"itemid"`
}

func main() {
	flag.Parse()

	if *updateFileFlag {
		fmt.Printf("Flag passed %t\n", *updateFileFlag)
		if err := getURLToFile(); err != nil {
			fmt.Printf("error occured with getURLToFile: %v", err)
			return
		}
	}

	if _, err := os.Stat(htmlFile); err != nil {
		fmt.Printf("%s does not exist, please run with the -u flag to initialize.\n", htmlFile)
		return
	}

	//Allows you to simulate local file as HTTP request
	//Below, "." local directory is chosen and transport is hooked up to colly
	transport := &http.Transport{}
	transport.RegisterProtocol("file", http.NewFileTransport(http.Dir(".")))

	collector := colly.NewCollector()
	collector.WithTransport(transport)

	var itemList []item

	//Works on <select name=itemlist first then <option
	collector.OnHTML("select[name=itemlist] option", func(h *colly.HTMLElement) {
		//This returns the value attribute of <option value=""
		urlString := h.Attr("value")

		idStringNumber, err := parseItemNumber(urlString)
		if err != nil {
			fmt.Printf("error parsing url: %v", err)
			return
		}

		newItem := item{
			//This returns the text inbetween <option></option>
			Name: h.Text,
			ID:   idStringNumber,
		}
		itemList = append(itemList, newItem)
	})

	collector.Visit(htmlFilePath)

	//JSONify the data I got.
	content, err := json.Marshal(itemList)
	if err != nil {
		fmt.Printf("error marshalling json: %v", err)
		return
	}

	os.WriteFile(jsonFile, content, 0644)
	fmt.Printf("Wrote item data to %s", htmlFile)
}

// Gets html page and locally stores it in working directory.
func getURLToFile() error {
	fmt.Printf("Beginning to get data from %s\n", URL)
	// Requesting HTML page
	resp, err := http.Get(URL)
	if err != nil {
		return fmt.Errorf("failed getting URL: %w", err)
	}
	defer resp.Body.Close()

	outputFile, err := os.Create(htmlFile)
	if err != nil {
		return fmt.Errorf("failed creating file: %w", err)
	}
	defer outputFile.Close()

	if _, err := io.Copy(outputFile, resp.Body); err != nil {
		return fmt.Errorf("failed copying file content: %w", err)
	}
	fmt.Printf("Completed, data written to %s\n", htmlFile)
	return nil
}

// Parses a url string to find the item number which is preceded by "itemid="
func parseItemNumber(urlString string) (string, error) {
	//Regex to find itemid=0000 where 0000 is any id number from urlString
	reg, err := regexp.Compile(`itemid=(\d*)`)
	if err != nil {
		return "", fmt.Errorf("error compiling regex: %w", err)
	}

	//returns itemid=0000 slice where [0] is itemid=0000 and [1] is 0000
	idStringNumber := reg.FindStringSubmatch(urlString)[1]

	return idStringNumber, nil

}
