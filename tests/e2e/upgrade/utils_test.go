// This file contains unit tests for the e2e package.
package upgrade

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestAppVersionsLess tests the AppVersion type's Less method with
// different version strings
func TestAppVersionsLess(t *testing.T) {
	var version AppVersion

	testCases := []struct {
		Name string
		Ver  string
		Exp  bool
	}{
		{
			Name: "higher - v10.0.1",
			Ver:  "v10.0.1",
			Exp:  false,
		},
		{
			Name: "lower - v9.1.0",
			Ver:  "v9.1.0",
			Exp:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			version = []string{tc.Ver, "v10.0.0"}
			require.Equal(t, version.Less(0, 1), tc.Exp, "expected: %v, got: %v", tc.Exp, version)
		})
	}
}

// TestAppVersionsSwap tests the AppVersion type's Swap method
func TestAppVersionsSwap(t *testing.T) {
	var version AppVersion
	value := "v9.1.0"
	version = []string{value, "v10.0.0"}
	version.Swap(0, 1)
	require.Equal(t, value, version[1], "expected: %v, got: %v", value, version[1])
}

// TestAppVersionsLen tests the AppVersion type's Len method
func TestAppVersionsLen(t *testing.T) {
	var version AppVersion = []string{"v9.1.0", "v10.0.0"}
	require.Equal(t, 2, version.Len(), "expected: %v, got: %v", 2, version.Len())
}

// TestRetrieveUpgradesList tests if the list of available upgrades in the codebase
// can be correctly retrieved
func TestRetrieveUpgradesList(t *testing.T) {
	upgradeList, err := RetrieveUpgradesList(upgradesPath)
	if os.IsNotExist(err) {
		t.Skip("skipping test as upgrade list file does not exist")
	}
	require.NoError(t, err, "expected no error while retrieving upgrade list")
	require.NotEmpty(t, upgradeList, "expected upgrade list to be non-empty")

	// check if all entries in the list match a semantic versioning pattern
	for _, upgrade := range upgradeList {
		require.Regexp(t, `^v\d+\.\d+\.\d+(-rc\d+)*$`, upgrade, "expected upgrade version to be in semantic versioning format")
	}
}
