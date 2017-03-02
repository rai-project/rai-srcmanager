package cmd

import (
	"github.com/k0kubun/pp"
	"github.com/spf13/cobra"
)

type Version struct {
	BuildDate  string
	GitCommit  string
	GitBranch  string
	GitState   string
	GitSummary string
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"info", "v"},
	Short:   "Prints a the version information of rai-srcmanager",
	Run: func(cmd *cobra.Command, args []string) {
		pp.Println(version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
