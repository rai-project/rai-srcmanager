package cmd

import (
	srcmanager "github.com/rai-project/rai-srcmanager/pkg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"clone", "up"},
	Short:   "Pull or clone rai repositories",
	Run: func(cmd *cobra.Command, args []string) {
		err := srcmanager.Update(isPublic)
		if err != nil {
			log.WithError(err).Fatal("Cannot update repositories")
		}
	},
}

func init() {
	updateCmd.PersistentFlags().BoolVar(&isPublic, "public", false, "use public repositories")
	RootCmd.AddCommand(updateCmd)
}
