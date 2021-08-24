package cmd

import(
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command {
	Use:	"",
	Short:  "",
	Example: "",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}