package config

import "fmt"

//Version info
var (
	GitCommit string
	GitTag    string
	BuildDate string

	Name        = "x-edge"
	Description = "A go-micro edge server app"
	version     = "2.0.0"
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
