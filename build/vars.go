package build

const Unspecified = "Unspecified"

// BuildID is set by the compiler (using -ldflags "-X util.BuildID $(git rev-parse --short HEAD)")
// and is used by GetID
var BuildID string

// BuildBranch is set by the compiler and is used by GetBranch
var BuildBranch string

// GetBuildID identifies what build is running.
func GetID() (BuildID string) {
	if BuildID != "" {
		return BuildID
	}

	return Unspecified
}

// GetBranch identifies the building host
func GetBranch() (BuildBranch string) {
	if BuildBranch != "" {
		return BuildBranch
	}

	return Unspecified
}
