package main

import (
	"runtime/debug"
	"testing"
)

func TestVersionString(t *testing.T) {
	tests := []struct {
		name    string
		buildFn readBuildInfoFunc
		got     string
		want    string
	}{
		{
			name: "VersionFromLocalBuild(Tagless)",
			buildFn: func() (info *debug.BuildInfo, ok bool) {
				// {"myBinVersion":"v1.0.7-0.20250421005157-8a88105c89cf"
				// "kubernetesVendor":"1.31"
				// "gitCommit":"8a88105c89cf"
				// "buildDate":"2025-04-21T00:51:57Z"
				// "goOs":"linux"
				// "goArch":"amd64"}
				return &debug.BuildInfo{
					Main: debug.Module{
						Version: "devel",
					},
					Settings: []debug.BuildSetting{
						{
							Key:   "GOOS",
							Value: "linux",
						},
						{
							Key:   "GOARCH",
							Value: "amd64",
						},
						{
							Key:   "vcs.revision",
							Value: "v1.0.7-0.20250421005157-8a88105c89cf",
						},
						{
							Key:   "vcs.time",
							Value: "2025-04-21T00:51:57Z",
						},
					},
				}, true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := versionString(tt.buildFn); got != tt.want {
				t.Errorf("versionString() = %v, want %v", got, tt.want)
			}
		})
	}
}
