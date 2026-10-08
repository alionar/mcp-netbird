package tools

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	mcpnetbird "github.com/aantti/mcp-netbird"
)

func TestListNetbirdGroupsDecodesResourceObjects(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/groups" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{
				"id": "group1",
				"name": "res-amf-uat-all",
				"peers": [],
				"peers_count": 1,
				"resources": [
					{
						"id": "resource1",
						"name": "amf-uat-trusted",
						"type": "network",
						"address": "10.101.0.0/17"
					}
				],
				"resources_count": 1
			}
		]`))
	}))
	defer server.Close()

	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(server.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()

	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	groups, err := listNetbirdGroups(ctx, ListNetbirdGroupsParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(groups) != 1 {
		t.Fatalf("got %d groups, want 1", len(groups))
	}
	if len(groups[0].Resources) != 1 {
		t.Fatalf("got %d resources, want 1", len(groups[0].Resources))
	}

	resource := groups[0].Resources[0]
	if resource.ID != "resource1" || resource.Address != "10.101.0.0/17" {
		t.Errorf("unexpected resource: %+v", resource)
	}
}
