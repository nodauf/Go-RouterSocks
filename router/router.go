package router

import (
	"fmt"
	utils "github.com/nodauf/Go-RouterSocks/utils"
	"net"
	"strings"
)

var Routes = map[string]string{}

func AddRoutes(destination string, remoteSocks string) {
	destination = strings.ToLower(destination)
	if _, ok := Routes[destination]; ok {
		fmt.Println("[-] Route already present")
	} else {
		Routes[destination] = remoteSocks
		fmt.Println("[*] Successfull route added for network " + destination)
	}
}

func DeleteRoutes(destination string) {
	destination = strings.ToLower(destination)
	if _, ok := Routes[destination]; ok {
		delete(Routes, destination)
		fmt.Println("[*] Successfull route " + destination + " deleted")
	} else {
		fmt.Println("[-] Route not found")
	}
}

func PrintRoutes() {
	for destination, remoteSocks := range Routes {
		fmt.Println(destination + " => " + remoteSocks)
	}
}

func FlushRoutes() {
	for destination, _ := range Routes {
		delete(Routes, destination)
	}

	fmt.Println("[*] Successfull route flushed")
}

func DumpRoutes() {
	for destination, remoteSocks := range Routes {
		fmt.Println("route add " + destination + " " + remoteSocks)
	}

	fmt.Println("[*] Successfull route dumped")
}

func GetRoute(ip string) (string, string) {
	if geoipEnable {
		ipNet := net.ParseIP(ip)
		record, err := geoIPDB.Country(ipNet)
		if err != nil {
			fmt.Printf("[-] Fail to retrieve the country of %s\n", ip)
		}
		isoCode := strings.ToLower(record.Country.IsoCode)
		if _, exist := Routes[isoCode]; exist {
			return isoCode, Routes[isoCode]
		}
	}
	for destination, remoteSocks := range Routes {
		// if destination is iso code
		if IsValidIsoCode(destination) {
			continue
			//for _, network := range isoWithNetworks[destination] {
			//	if utils.CIDRContainsIP(network, ip) {
			//		return destination, remoteSocks
			//	}
			//}

		} else if utils.CIDRContainsIP(destination, ip) {
			return destination, remoteSocks
		}
	}
	return "", ""
}
