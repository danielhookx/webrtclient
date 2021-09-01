package cmd

import (
	"fmt"
	"github.com/pion/stun"
	"github.com/spf13/cobra"
)

var stunCmd = &cobra.Command{
	Use:     "stun",
	Short:   "stun client",
	Long:    "",
	Example: "stun -h stun.l.google.com:19302",
	Version: "",
	Run:     DoStun,
}

var (
	stunServerHost string
)

func init() {
	stunCmd.Flags().StringVarP(&stunServerHost, "host", "H", "stun.l.google.com:19302", "stun server host address")

	rootCmd.AddCommand(stunCmd)
}

func DoStun(cmd *cobra.Command, args []string) {
	// Creating a "connection" to STUN server.
	c, err := stun.Dial("udp", stunServerHost)
	if err != nil {
		panic(err)
	}
	// Building binding request with random transaction id.
	message := stun.MustBuild(stun.TransactionID, stun.BindingRequest)
	// Sending request to STUN server, waiting for response message.
	if err := c.Do(message, func(res stun.Event) {
		if res.Error != nil {
			panic(res.Error)
		}
		// Decoding XOR-MAPPED-ADDRESS attribute from message.
		var xorAddr stun.XORMappedAddress
		if err := xorAddr.GetFrom(res.Message); err != nil {
			panic(err)
		}
		fmt.Println("your IP is", xorAddr.IP)
		fmt.Println("your Port is", xorAddr.Port)
	}); err != nil {
		panic(err)
	}
}
