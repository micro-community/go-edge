package config

import "fmt"

//Version info
var (
	GitCommit string
	GitTag    string
	BuildDate string

	name        = "x-edge"
	description = "A go-micro edge server"
	version     = "1.18.0"
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
