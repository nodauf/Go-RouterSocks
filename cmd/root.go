package cmd

import (
	"fmt"
	"github.com/nodauf/Go-RouterSocks/router"
	"log"
	"os"

	"github.com/spf13/cobra"

	prompt "github.com/nodauf/Go-RouterSocks/prompt"
	socks "github.com/nodauf/Go-RouterSocks/socks"
)

var cfgFile string
var port int
var ip string
var geoIPDB string

var rootCmd = &cobra.Command{
	Use:   "rsocks",
	Short: "Router socks",
	Long: `Run a socks server and redirect the traffic to other socks server
	according of the defined routes.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if geoIPDB != "" {
			if _, err := os.Open(geoIPDB); err != nil {
				return fmt.Errorf("cannot open the file %s: %s", geoIPDB, err.Error())
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if geoIPDB != "" {
			if err := router.LoadGeoIPDatabase(geoIPDB); err != nil {
				log.Fatalf("error while loading the GeoIP database: %s", err.Error())
			} else {
				log.Printf("[*] GeoIP database %s loaded\n", geoIPDB)
			}
		}
		socks.StartSocks(ip, port)
		prompt.Prompt()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	rootCmd.Flags().IntVarP(
		&port, "port", "p", 1080,
		"Socks5 port",
	)
	rootCmd.Flags().StringVarP(
		&ip, "ip", "i", "0.0.0.0",
		"IP for socks5 server",
	)
	rootCmd.Flags().StringVarP(
		&geoIPDB, "geoip", "g", "",
		"Path to the GeoIP database (GeoLite2-Country.mmdb)",
	)
}
