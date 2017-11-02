package cmd

import (
	"fmt"
	"os"

	srcmanager "github.com/rai-project/rai-srcmanager/pkg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	Verbose  bool
	isPublic bool
	version  Version
	// RootCmd represents the base command when called without any subcommands
	RootCmd = &cobra.Command{
		Use:   "rai-srcmanager",
		Short: "rai-srcmanager is the utility for managing the rai repositories",
	}
)

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(v Version) {
	version = v
	if Verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	log.SetFormatter(new(log.TextFormatter))
	srcmanager.Logger = log.StandardLogger()
	srcmanager.Verbose = Verbose
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", true, "verbose output")
}
