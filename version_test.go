package refyne

import "testing"

func TestParseVersion(t *testing.T) {
	tests := []struct {
		version    string
		major      int
		minor      int
		patch      int
		prerelease string
	}{
		{"1.2.3", 1, 2, 3, ""},
		{"0.0.0", 0, 0, 0, ""},
		{"10.20.30", 10, 20, 30, ""},
		{"1.2.3-beta", 1, 2, 3, "beta"},
		{"1.2.3-beta.1", 1, 2, 3, "beta.1"},
		{"invalid", 0, 0, 0, ""},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			major, minor, patch, prerelease := ParseVersion(tt.version)

			if major != tt.major {
				t.Errorf("major = %d, want %d", major, tt.major)
			}
			if minor != tt.minor {
				t.Errorf("minor = %d, want %d", minor, tt.minor)
			}
			if patch != tt.patch {
				t.Errorf("patch = %d, want %d", patch, tt.patch)
			}
			if prerelease != tt.prerelease {
				t.Errorf("prerelease = %q, want %q", prerelease, tt.prerelease)
			}
		})
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"1.2.3", "1.2.3", 0},
		{"2.0.0", "1.0.0", 1},
		{"1.0.0", "2.0.0", -1},
		{"1.2.0", "1.1.0", 1},
		{"1.1.0", "1.2.0", -1},
		{"1.1.2", "1.1.1", 1},
		{"1.1.1", "1.1.2", -1},
	}

	for _, tt := range tests {
		name := tt.a + " vs " + tt.b
		t.Run(name, func(t *testing.T) {
			got := CompareVersions(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("CompareVersions(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestVersionConstants(t *testing.T) {
	t.Run("SDK version is valid", func(t *testing.T) {
		major, _, _, _ := ParseVersion(SDKVersion)
		if major < 0 {
			t.Error("invalid SDK version")
		}
	})

	t.Run("min <= max", func(t *testing.T) {
		cmp := CompareVersions(MinAPIVersion, MaxKnownAPIVersion)
		if cmp > 0 {
			t.Error("MinAPIVersion > MaxKnownAPIVersion")
		}
	})
}
