package main

import "github.com/rai-project/rai-srcmanager/cmd"

var (
	// These fields are populated by govvv
	BuildDate  string
	GitCommit  string
	GitBranch  string
	GitState   string
	GitSummary string
	version    cmd.Version
)

func init() {
	version = cmd.Version{
		BuildDate:  BuildDate,
		GitCommit:  GitCommit,
		GitBranch:  GitBranch,
		GitState:   GitState,
		GitSummary: GitSummary,
	}
}
