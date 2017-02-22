package cmd

import (
	log "github.com/Sirupsen/logrus"
	srcmanager "github.com/rai-project/rai-srcmanager/pkg"
	"github.com/spf13/cobra"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit [msg]",
	Short: "Adds and records all the changes to the repository",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("Aborting commit due to empty commit message")
			return
		}
		srcmanager.Commit(args[0])
	},
}

func init() {
	RootCmd.AddCommand(commitCmd)
}
