package integration

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/AndreyShep2012/go-company-handler/internal/version"
	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	expected := version.VersionResponse{
		Name:      version.AppName,
		Version:   "0.0.0",
		BuildTime: "time",
		Revision:  "1.0",
		Branch:    "main",
	}

	version.Branch = expected.Branch
	version.BuildTime = expected.BuildTime
	version.Revision = expected.Revision
	version.Version = expected.Version

	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://"+testConf.ListenAddr+"/version", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var res version.VersionResponse
	err = json.Unmarshal(bodyBytes, &res)
	require.NoError(t, err)

	require.Equal(t, expected, res)
}
