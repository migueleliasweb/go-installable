package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// var needs to be used instead of const as ldflags is used to fill this
// information in the release process
var (
	myBinVersion            = ""
	kubernetesVendorVersion = "1.31"
	goos                    = ""
	goarch                  = ""
	gitCommit               = "" // "$Format:%H$" sha1 from git, output of $(git rev-parse HEAD)
	buildDate               = "" // "1970-01-01T00:00:00Z" build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
)

// version contains all the information related to the CLI version
type version struct {
	MyBinVersion     string `json:"myBinVersion,omitempty"`
	KubernetesVendor string `json:"kubernetesVendor,omitempty"`
	GitCommit        string `json:"gitCommit,omitempty"`
	BuildDate        string `json:"buildDate,omitempty"`
	GoOs             string `json:"goOs,omitempty"`
	GoArch           string `json:"goArch,omitempty"`
}

func buildInfoSettingToMap(settings []debug.BuildSetting) map[string]string {
	result := map[string]string{}

	for _, settingValue := range settings {
		result[settingValue.Key] = settingValue.Value
	}

	return result
}

func cleanupBuildDateFromMainVersionSplit(buildDateSplit string) string {
	buildDatesplit := strings.Split(
		buildDateSplit,
		".",
	)

	if len(buildDatesplit) == 2 {
		return buildDatesplit[1]
	}

	return buildDateSplit
}

// versionString returns the CLI version
func versionString() string {
	if myBinVersion == "" {
		if buildInfo, ok := debug.ReadBuildInfo(); ok {
			if buildInfo.Main.Version != "" {
				myBinVersion = buildInfo.Main.Version
			} else {
				myBinVersion = "(devel)"
			}

			buildSettingMap := buildInfoSettingToMap(buildInfo.Settings)

			goos = buildSettingMap["GOOS"]
			goarch = buildSettingMap["GOARCH"]

			gitCommitValue, gitCommitFound := buildSettingMap["vcs.revision"]

			if gitCommitFound {
				gitCommit = gitCommitValue
			}

			vcsTimeValue, vcsTimeFound := buildSettingMap["vcs.time"]

			if vcsTimeFound {
				buildDate = vcsTimeValue
			}

			// fallback to `.Main.Version`
			if !gitCommitFound && !vcsTimeFound && buildInfo.Main.Version != "" {
				// (usual)		<semver>-<build-date>-<commit-hash>
				// (sometimes) 	<semver>-<number>.<build-date>-<commit-hash>
				mainVersionSplit := strings.Split(buildInfo.Main.Version, "-")

				gitCommit = mainVersionSplit[2]

				if t, err := time.Parse("20060102150405", cleanupBuildDateFromMainVersionSplit(mainVersionSplit[1])); err == nil {
					buildDate = t.Format(time.RFC3339)
				}
			}
		}
	}

	jsonBytes, err := json.Marshal(version{
		MyBinVersion:     myBinVersion,
		KubernetesVendor: kubernetesVendorVersion,
		GitCommit:        gitCommit,
		BuildDate:        buildDate,
		GoOs:             goos,
		GoArch:           goarch,
	})

	if err != nil {
		fmt.Println("Err:", err)
		os.Exit(1)
	}

	return string(jsonBytes)
}

func main() {
	fmt.Println(versionString())

	buildInfo, _ := debug.ReadBuildInfo()
	fmt.Println("deps:", buildInfo.Deps)
	fmt.Println("goversion:", buildInfo.GoVersion)
	fmt.Println("main:", buildInfo.Main)
	fmt.Println("path:", buildInfo.Path)
	fmt.Println("settings:", buildInfo.Settings)
	fmt.Println("string:", buildInfo.String())
}
