package prompt

import "fmt"

func helpRoute() {
	fmt.Println("Route usage: ")
	helpRouteAdd()
	fmt.Print("\n")
	helpRouteDelete()
}

func helpRouteAdd() {
	fmt.Println(`route add cidr socksServer:socksPort
route add cidr chiselID
  route add 192.168.1.0/24 127.0.0.1:1081
  route add 192.168.1.0/24 0`)
}

func helpRouteDelete() {
	fmt.Println(`route delete cidr
  route delete 192.168.1.0/24`)
}

func helpChisel() {
	fmt.Println(`chisel
  Output:
  [0] 127.0.0.1:1081`)
}

func help() {
	fmt.Println("route: Manage route to socks servers")
	fmt.Println("chisel: Liste chisel socks server on localhost")
	fmt.Println("help: help command")
}
