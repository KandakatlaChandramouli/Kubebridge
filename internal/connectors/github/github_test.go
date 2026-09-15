package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KandakatlaChandramouli/kubebridge/internal/connectors"
)

func TestConnectorName(t *testing.T) {
	connector := New()

	assert.Equal(t, "github", connector.Name())
}

func TestConnectRequiresToken(t *testing.T) {
	connector := New()

	err := connector.Connect(context.Background(), connectors.Config{})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "github token is required")
}

func TestListResourcesWithoutConnection(t *testing.T) {
	connector := New()

	resources, err := connector.ListResources(context.Background())

	require.Error(t, err)
	assert.Nil(t, resources)
	assert.Contains(t, err.Error(), "not connected")
}

func TestListResourcesFromOrganization(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/orgs/kubernetes/repos", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(`[
			{
				"id": 1,
				"name": "kubernetes",
				"full_name": "kubernetes/kubernetes",
				"html_url": "https://github.com/kubernetes/kubernetes",
				"language": "Go"
			},
			{
				"id": 2,
				"name": "website",
				"full_name": "kubernetes/website",
				"html_url": "https://github.com/kubernetes/website",
				"language": "JavaScript"
			}
		]`))

		require.NoError(t, err)
	}))
	defer server.Close()

	connector := New()

	err := connector.Connect(
		context.Background(),
		connectors.Config{
			"token":        "test-token",
			"organization": "kubernetes",
		},
	)
	require.NoError(t, err)

	connector.client.BaseURL, _ = connector.client.BaseURL.Parse(server.URL + "/")
	connector.client.UploadURL, _ = connector.client.UploadURL.Parse(server.URL + "/")

	resources, err := connector.ListResources(context.Background())

	require.NoError(t, err)
	require.Len(t, resources, 2)

	assert.Equal(t, "kubernetes/kubernetes", resources[0].ID)
	assert.Equal(t, "kubernetes", resources[0].Name)
	assert.Equal(t, "github_repository", resources[0].Type)
	assert.Equal(
		t,
		"https://github.com/kubernetes/kubernetes",
		resources[0].Metadata["url"],
	)
	assert.Equal(t, "Go", resources[0].Metadata["language"])

	assert.Equal(t, "kubernetes/website", resources[1].ID)
	assert.Equal(t, "website", resources[1].Name)
}

func TestDisconnect(t *testing.T) {
	connector := New()

	err := connector.Connect(
		context.Background(),
		connectors.Config{"token": "test-token"},
	)
	require.NoError(t, err)

	err = connector.Disconnect(context.Background())
	require.NoError(t, err)

	resources, err := connector.ListResources(context.Background())

	require.Error(t, err)
	assert.Nil(t, resources)
}
