package refyne

import (
	"fmt"
	"regexp"
	"strconv"
)

var versionRegex = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)(?:-(.+))?$`)

// ParseVersion parses a semver version string into components.
func ParseVersion(version string) (major, minor, patch int, prerelease string) {
	match := versionRegex.FindStringSubmatch(version)
	if match == nil {
		return 0, 0, 0, ""
	}

	major, _ = strconv.Atoi(match[1])
	minor, _ = strconv.Atoi(match[2])
	patch, _ = strconv.Atoi(match[3])
	if len(match) > 4 {
		prerelease = match[4]
	}

	return
}

// CompareVersions compares two semver versions.
// Returns -1 if a < b, 0 if a == b, 1 if a > b.
func CompareVersions(a, b string) int {
	aMajor, aMinor, aPatch, _ := ParseVersion(a)
	bMajor, bMinor, bPatch, _ := ParseVersion(b)

	if aMajor != bMajor {
		if aMajor < bMajor {
			return -1
		}
		return 1
	}
	if aMinor != bMinor {
		if aMinor < bMinor {
			return -1
		}
		return 1
	}
	if aPatch != bPatch {
		if aPatch < bPatch {
			return -1
		}
		return 1
	}

	return 0
}

// CheckAPIVersionCompatibility checks if an API version is compatible with this SDK.
func CheckAPIVersionCompatibility(apiVersion string, logger Logger) error {
	// If API version is lower than minimum supported, return error
	if CompareVersions(apiVersion, MinAPIVersion) < 0 {
		return &UnsupportedAPIVersionError{
			APIVersion:      apiVersion,
			MinVersion:      MinAPIVersion,
			MaxKnownVersion: MaxKnownAPIVersion,
		}
	}

	// If API major version is higher than known, warn
	apiMajor, _, _, _ := ParseVersion(apiVersion)
	maxMajor, _, _, _ := ParseVersion(MaxKnownAPIVersion)

	if apiMajor > maxMajor {
		logger.Warn(
			fmt.Sprintf(
				"API version %s is newer than this SDK was built for (%s). "+
					"There may be breaking changes. Consider upgrading the SDK.",
				apiVersion, MaxKnownAPIVersion,
			),
			map[string]any{
				"apiVersion":      apiVersion,
				"sdkVersion":      SDKVersion,
				"maxKnownVersion": MaxKnownAPIVersion,
			},
		)
	}

	return nil
}
