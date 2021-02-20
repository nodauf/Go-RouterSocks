package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	prompt "github.com/nodauf/Go-RouterSocks/prompt"
	socks "github.com/nodauf/Go-RouterSocks/socks"
)

var cfgFile string
var port int
var ip string

var rootCmd = &cobra.Command{
	Use:   "rsocks",
	Short: "Router socks",
	Long: `Run a socks server and redirect the traffic to other socks server
	according of the defined routes.`,
	Run: func(cmd *cobra.Command, args []string) {
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
}
