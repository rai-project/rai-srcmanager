package cmd

import (
	log "github.com/Sirupsen/logrus"
	srcmanager "github.com/rai-project/rai-srcmanager/pkg"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Pull or clone rai repositories",
	Run: func(cmd *cobra.Command, args []string) {
		err := srcmanager.Update()
		if err != nil {
			log.WithError(err).Fatal("Cannot update repositories")
		}
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)
}
