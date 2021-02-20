package router

import (
	"fmt"

	utils "github.com/nodauf/Go-RouterSocks/utils"
)

var Routes = map[string]string{}

func AddRoutes(network string, remoteSocks string) {
	if _, ok := Routes[network]; ok {
		fmt.Println("[-] Route already present")
	} else {
		Routes[network] = remoteSocks
		fmt.Println("[*] Successfull route added")
	}
}

func DeleteRoutes(network string) {
	if _, ok := Routes[network]; ok {
		delete(Routes, network)
		fmt.Println("[*] Successfull route deleted")
	} else {
		fmt.Println("[-] Route not found")
	}
}

func PrintRoutes() {
	for network, remoteSocks := range Routes {
		fmt.Println(network + " => " + remoteSocks)
	}
}

func FlushRoutes() {
	for network, _ := range Routes {
		delete(Routes, network)
	}

	fmt.Println("[*] Successfull route flushed")
}

func GetRoute(ip string) string {
	for network, remoteSocks := range Routes {
		if utils.CIDRContainsIP(network, ip) {
			return remoteSocks
		}
	}
	return ""
}
