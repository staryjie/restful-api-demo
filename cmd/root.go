package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/staryjie/restful-api-demo/version"
)

var vers bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "demo-api",
	Short: "demo-api 后端API",
	Long:  "demo-api 后端API",
	RunE: func(cmd *cobra.Command, args []string) error {
		if vers {
			fmt.Println(version.FullVersion())
			return nil
		}
		return nil
	},
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&vers, "version", "v", false, "print demo-api version")
}
