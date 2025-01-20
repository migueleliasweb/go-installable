package main

import (
	"fmt"
	"runtime/debug"
)

const unknown = "unknown"

// var needs to be used instead of const as ldflags is used to fill this
// information in the release process
var (
	myBinVersion            = unknown
	kubernetesVendorVersion = "1.31"
	goos                    = unknown
	goarch                  = unknown
	gitCommit               = "$Format:%H$" // sha1 from git, output of $(git rev-parse HEAD)

	buildDate = "1970-01-01T00:00:00Z" // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
)

// version contains all the information related to the CLI version
type version struct {
	MyBinVersion     string `json:"myBinVersion"`
	KubernetesVendor string `json:"kubernetesVendor"`
	GitCommit        string `json:"gitCommit"`
	BuildDate        string `json:"buildDate"`
	GoOs             string `json:"goOs"`
	GoArch           string `json:"goArch"`
}

func buildInfoSettingTpMap(settings []debug.BuildSetting) map[string]string {
	result := map[string]string{}

	for _, settingValue := range settings {
		result[settingValue.Key] = settingValue.Value
	}

	return result
}

// versionString returns the CLI version
func versionString() string {
	if myBinVersion == unknown {
		if buildInfo, ok := debug.ReadBuildInfo(); ok {
			if buildInfo.Main.Version != "" {
				myBinVersion = buildInfo.Main.Version
			}

			buildSettingMap := buildInfoSettingTpMap(buildInfo.Settings)

			goos = buildSettingMap["GOOS"]
			goarch = buildSettingMap["GOARCH"]

			if v, found := buildSettingMap["vcs.revision"]; found {
				gitCommit = v
			}
		}
	}

	return fmt.Sprintf("Version: %#v", version{
		myBinVersion,
		kubernetesVendorVersion,
		gitCommit,
		buildDate,
		goos,
		goarch,
	})
}

func main() {
	fmt.Println(versionString())
}
