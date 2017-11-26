package cmd

import (
	srcmanager "github.com/rai-project/rai-srcmanager/pkg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	ignoredCommitFlags string
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:     "commit [msg]",
	Aliases: []string{"co"},
	Short:   "Adds and records all the changes to the repository",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("Aborting commit due to empty commit message")
			return
		}
		srcmanager.Commit(isPublic, args[0])
	},
}

func init() {
	commitCmd.PersistentFlags().StringVarP(&ignoredCommitFlags, "add", "a", "", "add all this is the default and cannot be overridden")
	commitCmd.PersistentFlags().StringVarP(&ignoredCommitFlags, "message", "m", "", "message this is the default and cannot be overridden")
	commitCmd.PersistentFlags().StringVar(&ignoredCommitFlags, "am", "", "add all with message this is the default and cannot be overridden")
	commitCmd.PersistentFlags().BoolVar(&isPublic, "public", false, "use public repositories")

	RootCmd.PersistentFlags().MarkHidden("add")
	RootCmd.PersistentFlags().MarkHidden("message")
	RootCmd.PersistentFlags().MarkHidden("am")

	RootCmd.AddCommand(commitCmd)
}
