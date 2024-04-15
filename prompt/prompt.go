package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	router "github.com/nodauf/Go-RouterSocks/router"
	utils "github.com/nodauf/Go-RouterSocks/utils"

	prompt "github.com/c-bata/go-prompt"
)

func executor(in string) {
	in = strings.TrimSpace(in)
	command := strings.Split(in, " ")
	first := command[0]
	switch strings.ToLower(first) {
	case "help":
		if len(command) > 1 {
			second := command[1]
			switch strings.ToLower(second) {
			case "route":
				helpRoute()
			case "chisel":
				helpChisel()
			case "geoip":
				helpGeoIP()
			}
		} else {
			help()
		}
	case "route":
		if len(command) > 1 {
			second := command[1]
			switch strings.ToLower(second) {
			case "list":
				router.PrintRoutes()
			case "add":
				if len(command) != 4 {
					helpRouteAdd()
				} else if serverSocks := utils.IsChiselIDValid(command[3]); serverSocks != "" {
					remoteNetwork := command[2]
					router.AddRoutes(remoteNetwork, serverSocks)
				} else if !utils.IsCIDRValid(command[2]) && !router.IsValidIsoCode(command[2]) {
					fmt.Println("[-] CIDR or ISO code is not valid")
					helpRouteAdd()
				} else if !utils.IsRemoteSocksValid(command[3]) {
					fmt.Println("[-] Socks server, socks port or chisel ID is not valid")
					helpRouteAdd()
				} else if !utils.CanResolvedHostname(strings.Split(command[3], ":")[0]) {
					fmt.Println("[-] Server socks can be resolved")
					helpRouteAdd()
				} else if !utils.ServerReachable(command[3]) {
					fmt.Println("[-] Server is not reachable")
					helpRouteAdd()
				} else {
					remoteNetwork := command[2]
					remoteSocks := command[3]
					router.AddRoutes(remoteNetwork, remoteSocks)
				}

			case "delete":
				if len(command) != 3 {
					helpRouteDelete()
				} else if !utils.IsCIDRValid(command[2]) {
					fmt.Println("[-] CIDR is not valid")
					helpRouteDelete()
				} else {
					router.DeleteRoutes(command[2])
				}

			case "flush":
				router.FlushRoutes()

			case "dump":
				router.DumpRoutes()

			case "import":
				fmt.Println("Paste the output of a route dump and end with an empty line")
				var lines string
				scanner := bufio.NewScanner(os.Stdin)
				for scanner.Scan() {
					line := scanner.Text()
					lines += line + "\n"
					if line == "" {
						break
					}
				}
				for _, line := range strings.Split(lines, "\n") {
					if line != "" {
						executor(line)
					}
				}
			default:

				fmt.Println("Invalid route command")
			}
		} else {
			router.PrintRoutes()
		}
	case "chisel":
		utils.PrintChiselProcess()
	case "geoip":
		if len(command) > 1 {
			second := command[1]
			switch strings.ToLower(second) {
			case "load":

				if len(command) != 3 {
					helpGeoIPLoad()
				} else {
					err := router.LoadGeoIPDatabase(command[2])
					if err != nil {
						fmt.Printf("[-] Fail to load the database %s: %s\n", command[2], err.Error())
					} else {
						fmt.Printf("[*] GeoIP database %s loaded\n", command[2])
					}
				}
			case "print":
				router.PrintCountry()
			}
		} else {
			helpGeoIP()
		}
	case "exit":
		os.Exit(0)
	case "":
	default:
		fmt.Println("Invalid command")
	}
}

func Prompt() {
	p := prompt.New(
		executor,
		complete,
		prompt.OptionPrefix("RouterSocks> "),
		prompt.OptionPrefixTextColor(prompt.Red),
		prompt.OptionTitle("Router Socks"),
	)
	p.Run()
	//fmt.Println("You selected " + t)
}
