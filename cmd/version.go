package cmd

import "fmt"

//Version info
var (
	GitCommit string
	GitTag    string
	BuildDate string

	Name        = "x-edge"
	Description = "an edge connection framework"
	version     = "2.0.0-alpha"
)

//BuildVersion for framework
func BuildVersion() string {
	microVersion := version

	if GitTag != "" {
		microVersion = GitTag
	}

	if GitCommit != "" {
		microVersion += fmt.Sprintf("-%s", GitCommit)
	}

	if BuildDate != "" {
		microVersion += fmt.Sprintf("-%s", BuildDate)
	}

	return microVersion
}
