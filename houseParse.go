package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type HouseLobbyist struct {
	FirstName string `xml:"lobbyistFirstName"`
	LastName  string `xml:"lobbyistLastName"`
}

type HouseFiling struct {
	OrganizationName string          `xml:"organizationName"`
	ClientName       string          `xml:"clientName"`
	SenateID         string          `xml:"senateID"`
	HouseID          string          `xml:"houseID"`
	ReportYear       string          `xml:"reportYear"`
	ReportType       string          `xml:"reportType"`
	Lobbyist         []HouseLobbyist `xml:"alis>ali_info>lobbyists>lobbyist"` //different formats for quarterly vs aggregate reports?
	//Lobbyist []Lobbyist `xml:"lobbyists>lobbyist"`
}

func parseHouseFilings(recordDir string, combinedFilings *[]GenericFiling, mutex *sync.Mutex, wg *sync.WaitGroup) {
	beginParseTime := time.Now()

	files, err := ioutil.ReadDir("./" + recordDir + "/")
	if err != nil {
		panic(err)
	}

	fmt.Println("Reading " + strconv.Itoa(len(files)) + " files from " + recordDir + "...")

	a := 0 //counter for number of files successfully read

	for _, f := range files {
		data, err := ioutil.ReadFile(recordDir + "/" + f.Name())
		if err != nil {
			fmt.Println("error reading", err)
			continue
		} else {
			if strings.Contains(filepath.Ext(f.Name()), "xml") {
				oneFiling := HouseFiling{}
				//unmarshal data and put into struct array
				err = xml.Unmarshal(data, &oneFiling)
				if err != nil {
					fmt.Println("error decoding %v: %v", f.Name(), err)
					continue
				} else {
					mutex.Lock()
					combineSingleFiling(oneFiling, combinedFilings)
					mutex.Unlock()
					a++ //increment number of files successfully parsed
				}

			}
		}

		if a%10000 == 0 {
			fmt.Println(strconv.Itoa(a), "House files read")
		}
	}

	fmt.Println("Successfully read ", a, " / ", len(files), "House files in", time.Since(beginParseTime).String())

	fmt.Println("Removing record directory " + recordDir + "...")
	err = os.RemoveAll(recordDir)
	if err != nil {
		panic(err)
	}
	fmt.Println("Removed record directory " + recordDir)

	//Waitgroup done
	wg.Done()
}
