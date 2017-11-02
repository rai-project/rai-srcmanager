package cmd

import (
	srcmanager "github.com/rai-project/rai-srcmanager/pkg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// dirtyCmd represents the dirty command
var dirtyCmd = &cobra.Command{
	Use:     "dirty",
	Aliases: []string{"status"},
	Short:   "Prints a list of repositories with uncommitted changes",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			log.Warning("Ignoring args after 'status'")
		}
		srcmanager.Dirty(isPublic)
	},
}

func init() {
	dirtyCmd.PersistentFlags().BoolVar(&isPublic, "public", false, "use public repositories")
	RootCmd.AddCommand(dirtyCmd)
}
