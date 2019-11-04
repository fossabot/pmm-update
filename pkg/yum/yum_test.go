// pmm-update
// Copyright (C) 2019 Percona LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package yum

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstalled(t *testing.T) {
	res, err := Installed(context.Background(), "pmm-update")
	require.NoError(t, err)

	assert.True(t, strings.HasPrefix(res.Installed.Version, "2."), "%s", res.Installed.Version)
	assert.True(t, strings.HasPrefix(res.Installed.FullVersion, "2."), "%s", res.Installed.FullVersion)
	require.NotEmpty(t, res.Installed.BuildTime)
	assert.True(t, time.Since(*res.Installed.BuildTime) < 60*24*time.Hour, "InstalledTime = %s", res.Installed.BuildTime)
	assert.Equal(t, "local", res.Installed.Repo)
}

func TestCheck(t *testing.T) {
	res, err := Check(context.Background(), "pmm-update")
	require.NoError(t, err)

	assert.True(t, strings.HasPrefix(res.Installed.Version, "2."), "%s", res.Installed.Version)
	assert.True(t, strings.HasPrefix(res.Installed.FullVersion, "2."), "%s", res.Installed.FullVersion)
	require.NotEmpty(t, res.Installed.BuildTime)
	assert.True(t, time.Since(*res.Installed.BuildTime) < 60*24*time.Hour, "InstalledTime = %s", res.Installed.BuildTime)
	assert.Equal(t, "local", res.Installed.Repo)

	assert.True(t, strings.HasPrefix(res.Latest.Version, "2."), "%s", res.Latest.Version)
	assert.True(t, strings.HasPrefix(res.Latest.FullVersion, "2."), "%s", res.Latest.FullVersion)
	require.NotEmpty(t, res.Latest.BuildTime)
	assert.True(t, time.Since(*res.Latest.BuildTime) < 60*24*time.Hour, "LatestTime = %s", res.Latest.BuildTime)
	assert.NotEmpty(t, res.Latest.Repo)

	// We assume that the latest percona/pmm-server:2 and perconalab/pmm-server:dev-latest images
	// always contains the latest pmm-update package versions.
	// If this test fails, re-pull them and recreate devcontainer.
	var updateAvailable bool
	image := os.Getenv("PMM_SERVER_IMAGE")
	require.NotEmpty(t, image)
	if image != "percona/pmm-server:2" && image != "perconalab/pmm-server:dev-latest" {
		updateAvailable = true
	}
	if updateAvailable {
		t.Log("Assuming pmm-update update is available.")
		assert.True(t, res.UpdateAvailable, "update should be available")
		assert.True(t, strings.HasPrefix(res.LatestNewsURL, "https://per.co.na/pmm/2."), "latest_news_url = %q", res.LatestNewsURL)
		assert.NotEqual(t, res.Installed.Version, res.Latest.Version, "versions should not be the same")
		assert.NotEqual(t, res.Installed.FullVersion, res.Latest.FullVersion, "versions should not be the same")
		assert.NotEqual(t, *res.Installed.BuildTime, *res.Latest.BuildTime, "build times should not be the same (%s)", *res.Installed.BuildTime)
		assert.Equal(t, "pmm2-server", res.Latest.Repo)
	} else {
		t.Log("Assuming the latest pmm-update version.")
		assert.False(t, res.UpdateAvailable, "update should not be available")
		assert.Empty(t, res.LatestNewsURL, "latest_news_url should be empty")
		assert.Equal(t, res.Installed, res.Latest, "version should be the same (latest)")
		assert.Equal(t, *res.Installed.BuildTime, *res.Latest.BuildTime, "build times should be the same")
		assert.Equal(t, "local", res.Latest.Repo)
	}
}

func TestUpdate(t *testing.T) {
	err := Update(context.Background(), "make")
	require.NoError(t, err)
}