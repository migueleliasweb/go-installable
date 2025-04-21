package main

import (
	"fmt"
	"runtime/debug"
)

func main() {
	ogDebugInfo, ok := debug.ReadBuildInfo()

	fmt.Println("Original build info:", ogDebugInfo)
	fmt.Println("Original ok:", ok)

	// fmt.Println(versionString(debug.ReadBuildInfo))
}
