package utils

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	netstat "github.com/nodauf/Go-RouterSocks/utils/netstat"
)

var ChiselProcess = map[int]*netstat.SockAddr{}

func IsCIDRValid(cidr string) bool {
	var re = regexp.MustCompile(`^([0-9]{1,3}\.){3}[0-9]{1,3}\/([0-9]|[1-2][0-9]|3[0-2])$`)
	return re.Match([]byte(cidr))
}

func IsRemoteSocksValid(remoteSocks string) bool {
	var re = regexp.MustCompile(`^(\w+|(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-4]|2[0-4][0-9]|[01]?[0-9][0-9]?)):([1-9]|[1-5]?[0-9]{2,4}|6[1-4][0-9]{3}|65[1-4][0-9]{2}|655[1-2][0-9]|6553[1-5])$`)
	return re.Match([]byte(remoteSocks))
}

func CanResolvedHostname(server string) bool {
	_, err := net.LookupIP(server)
	if err != nil {
		return false
	}
	return true
}

func CIDRContainsIP(cidr string, ip string) bool {
	_, network, _ := net.ParseCIDR(cidr)
	fmt.Println(cidr)
	return network.Contains(net.ParseIP(ip))
}

func PrintChiselProcess() {
	GetChiselProcess()
	for id, server := range ChiselProcess {
		fmt.Println("[" + strconv.Itoa(id) + "] " + server.String())
	}
}

func IsChiselIDValid(idString string) string {
	// Init the chisel list in case the user didn't print it
	GetChiselProcess()
	id, err := strconv.Atoi(idString)
	if serverSocks, ok := ChiselProcess[id]; err == nil && ok {
		return serverSocks.String()
	}
	return ""
}

func GetChiselProcess() {
	ChiselProcess = map[int]*netstat.SockAddr{}
	tabs, err := netstat.TCPSocks(func(s *netstat.SockTabEntry) bool {
		return s.State == netstat.Listen && s.Process != nil && (strings.ToLower(s.Process.Name) == "chisel" || strings.ToLower(s.Process.Name) == "chisel.exe")
	})
	if err != nil {
		fmt.Println("Can list the listen port: " + err.Error())
	}
	if len(tabs) > 0 {
		for i, e := range tabs {
			ChiselProcess[i] = e.LocalAddr

		}
	}
}

func ServerReachable(server string) bool {
	timeout := 1 * time.Second
	_, err := net.DialTimeout("tcp", server, timeout)
	return err == nil
}
