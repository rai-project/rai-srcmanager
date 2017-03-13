package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// licenseCmd represents the license command
var licenseCmd = &cobra.Command{
	Use:   "license",
	Short: "Displays the project license.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(_escFSMustString(false, "/LICENSE.TXT"))
	},
}

func initLicense() {
	RootCmd.AddCommand(licenseCmd)
}
