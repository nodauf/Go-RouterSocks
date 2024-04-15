package prompt

import "fmt"

func helpRoute() {
	fmt.Println("Route usage: ")
	helpRouteAdd()
	fmt.Print("\n")
	helpRouteDelete()
}

func helpRouteAdd() {
	fmt.Println(`route add cidr/geoIPCountry socksServer:socksPort
route add cidr/geoIPCountry chiselID
  route add 192.168.1.0/24 127.0.0.1:1081
  route add 192.168.1.0/24 0`)
}

func helpRouteDelete() {
	fmt.Println(`route delete cidr/geoIPCountry
  route delete 192.168.1.0/24`)
}

func helpChisel() {
	fmt.Println(`chisel
  Output:
  [0] 127.0.0.1:1081`)
}

func helpGeoIP() {
	fmt.Println("geoip usage: ")
	helpGeoIPLoad()
	fmt.Print("\n")
	helpGeoPrint()
}

func helpGeoIPLoad() {
	fmt.Println(`geoip load Path_to_the_database
	geoip load /tmp/GeoLite2-Country.mmdb`)
}
func helpGeoPrint() {
	fmt.Println(`geoip print
  Output:
  iso code => Country name`)
}

func help() {
	fmt.Println("route: Manage route to socks servers")
	fmt.Println("chisel: Liste chisel socks server on localhost")
	fmt.Println("geoip: List the supported GeoIP country")
	fmt.Println("help: help command")
}
