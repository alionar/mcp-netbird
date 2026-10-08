package tools

import (
	"context"

	mcpnetbird "github.com/aantti/mcp-netbird"
	"github.com/mark3labs/mcp-go/server"
)

type NetbirdGroupMember struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type NetbirdGroupResource struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Address string `json:"address"`
}

type NetbirdGroup struct {
	ID             string                 `json:"id"`
	Issued         string                 `json:"issued"`
	Name           string                 `json:"name"`
	Peers          []NetbirdGroupMember   `json:"peers"`
	PeersCount     int                    `json:"peers_count"`
	Resources      []NetbirdGroupResource `json:"resources"`
	ResourcesCount int                    `json:"resources_count"`
}

type ListNetbirdGroupsParams struct{}

func listNetbirdGroups(ctx context.Context, args ListNetbirdGroupsParams) ([]NetbirdGroup, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient()
	}

	var groups []NetbirdGroup
	if err := client.Get(ctx, "/groups", &groups); err != nil {
		return nil, err
	}

	return groups, nil
}

var ListNetbirdGroups = mcpnetbird.MustTool(
	"list_netbird_groups",
	"List all Netbird groups",
	listNetbirdGroups,
)

func AddNetbirdGroupTools(mcp *server.MCPServer) {
	ListNetbirdGroups.Register(mcp)
}
