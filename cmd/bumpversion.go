package cmd

import (
	srcmanager "github.com/rai-project/rai-srcmanager/pkg"
	"github.com/spf13/cobra"
)

// gogetCmd represents the goget command
var bumpversionCmd = &cobra.Command{
	Use:   "bumpversion",
	Short: "bumpversion patch --commit && git push && git push --tags",
	Run: func(cmd *cobra.Command, args []string) {
		srcmanager.BumpVersion()
	},
}

func init() {
	RootCmd.AddCommand(bumpversionCmd)
}
