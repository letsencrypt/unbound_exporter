package util

import (
	"expvar"
)

const Unspecified = "Unspecified"

// BuildID is set by the compiler (using -ldflags "-X util.BuildID $(git rev-parse --short HEAD)")
// and is used by GetBuildID
var BuildID string

// BuildHost is set by the compiler and is used by GetBuildHost
var BuildHost string

// BuildTime is set by the compiler and is used by GetBuildTime
var BuildTime string

// BuildBranch is set by the compiler and is used by GetBuildBranch
var BuildBranch string

func init() {
	expvar.NewString("BuildID").Set(BuildID)
	expvar.NewString("BuildTime").Set(BuildTime)
}

// GetBuildID identifies what build is running.
func GetBuildID() (retID string) {
	retID = BuildID
	if retID == "" {
		retID = Unspecified
	}
	return
}

// GetBuildTime identifies when this build was made
func GetBuildTime() (retID string) {
	retID = BuildTime
	if retID == "" {
		retID = Unspecified
	}
	return
}

// GetBuildHost identifies the building host
func GetBuildHost() (retID string) {
	retID = BuildHost
	if retID == "" {
		retID = Unspecified
	}
	return
}

// GetBuildBranch identifies the building host
func GetBuildBranch() (retID string) {
	retID = BuildBranch
	if retID == "" {
		retID = Unspecified
	}
	return
}
