package router

import (
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"github.com/oschwald/maxminddb-golang"
	"log"
	"slices"
	"strings"
)

var geoIPDB *geoip2.Reader
var geoIPDB_low *maxminddb.Reader
var geoipEnable bool
var isoCodes []string
var isoWithCountries = make(map[string]string)

//var isoWithNetworks = make(map[string][]string)

func LoadGeoIPDatabase(geoipDB string) error {
	var err error
	geoIPDB, err = geoip2.Open(geoipDB)
	if err != nil {
		return err
	}
	geoIPDB_low, err = maxminddb.Open(geoipDB)
	geoipEnable = true
	// Load the countries and the iso code
	var record struct {
		Country struct {
			ISOCode      string            `maxminddb:"iso_code"`
			CountryNames map[string]string `maxminddb:"names"`
		} `maxminddb:"country"`
	} // Or any appropriate struct
	networks := geoIPDB_low.Networks(maxminddb.SkipAliasedNetworks)
	for networks.Next() {
		_, err := networks.Network(&record)
		if err != nil {
			log.Fatal(err)
		}
		ISOCode := strings.ToLower(record.Country.ISOCode)
		if !slices.Contains(isoCodes, ISOCode) {
			isoCodes = append(isoCodes, ISOCode)
		}
		isoWithCountries[ISOCode] = record.Country.CountryNames["en"]
		//isoWithNetworks[ISOCode] = append(isoWithNetworks[ISOCode], subnet.String())
	}
	if networks.Err() != nil {
		log.Panic(networks.Err())
	}
	slices.Sort(isoCodes)
	return err
}

func PrintCountry() {
	if !geoipEnable {
		fmt.Println("GeoIP database was not loaded")
		return
	}

	for _, isoCode := range isoCodes {
		fmt.Printf("%s => %s\n", isoCode, isoWithCountries[isoCode])
	}
}

func IsValidIsoCode(ISOCode string) bool {
	return slices.Contains(isoCodes, strings.ToLower(ISOCode))
}
