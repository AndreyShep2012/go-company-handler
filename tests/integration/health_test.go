package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHealth(t *testing.T) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://"+testConf.ListenAddr+"/health", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}
