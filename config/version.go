package config


var (
	GitCommit string
	GitTag    string
	BuildDate string

	name        = "x-edge"
	description = "A go-micro edge server"
	version     = "1.18.0"
)


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