// Copyright The KCL Authors. All rights reserved.

package version

// version will be set by build flags.
var version string

// GetVersionString() will return the latest version of kpm.
func GetVersionString() string {
	if len(version) == 0 {
		// If version is not set by build flags, return the version constant.
		return VersionTypeLatest.String()
	}
	return version
}

// VersionType is the version type of kpm.
type VersionType string

// String() will transform VersionType to string.
func (kvt VersionType) String() string {
	return string(kvt)
}

// All the kpm versions.
const (
	VersionTypeLatest = Version_0_7_0_alpha_1

	Version_0_7_0         VersionType = "0.7.0"
	Version_0_7_0_alpha_1 VersionType = "0.7.0-alpha.1"
)