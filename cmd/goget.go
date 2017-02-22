package cmd

import (
	srcmanager "github.com/rai-project/rai-srcmanager/pkg"
	"github.com/spf13/cobra"
)

// gogetCmd represents the goget command
var gogetCmd = &cobra.Command{
	Use:   "goget",
	Short: "Performs a glide install on all the packages that are needed by the repositories",
	Run: func(cmd *cobra.Command, args []string) {
		srcmanager.GoGet()
	},
}

func init() {
	RootCmd.AddCommand(gogetCmd)
}
