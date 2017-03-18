package main

import "github.com/rai-project/rai-srcmanager/cmd"

var (
	// These fields are populated by govvv
	Version    string
	BuildDate  string
	GitCommit  string
	GitBranch  string
	GitState   string
	GitSummary string
	version    cmd.Version
)

func init() {
	version = cmd.Version{
		Version:    Version,
		BuildDate:  BuildDate,
		GitCommit:  GitCommit,
		GitBranch:  GitBranch,
		GitState:   GitState,
		GitSummary: GitSummary,
	}
}
