package tools

import (
	"context"
	"fmt"
	"strings"

	mcpnetbird "github.com/aantti/mcp-netbird"
	"github.com/mark3labs/mcp-go/server"
)

type NetbirdPolicyRule struct {
	Action        string             `json:"action"`
	Bidirectional bool               `json:"bidirectional"`
	Description   string             `json:"description"`
	Destinations  []NetbirdPeerGroup `json:"destinations"`
	Enabled       bool               `json:"enabled"`
	ID            string             `json:"id"`
	Name          string             `json:"name"`
	Protocol      string             `json:"protocol"`
	Sources       []NetbirdPeerGroup `json:"sources"`
}

type NetbirdPolicy struct {
	Description         string              `json:"description"`
	Enabled             bool                `json:"enabled"`
	ID                  string              `json:"id"`
	Name                string              `json:"name"`
	Rules               []NetbirdPolicyRule `json:"rules"`
	SourcePostureChecks any                 `json:"source_posture_checks"`
}

type ListNetbirdPoliciesParams struct{}

func listNetbirdPolicies(ctx context.Context, args ListNetbirdPoliciesParams) ([]NetbirdPolicy, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient()
	}

	var policies []NetbirdPolicy
	if err := client.Get(ctx, "/policies", &policies); err != nil {
		return nil, err
	}

	return policies, nil
}

var ListNetbirdPolicies = mcpnetbird.MustTool(
	"list_netbird_policies",
	"List all Netbird policies",
	listNetbirdPolicies,
)

type FindNetbirdPoliciesParams struct {
	Group string `mcp:"group" validate:"required" json:"group"`
}

func findNetbirdPolicies(ctx context.Context, args FindNetbirdPoliciesParams) ([]NetbirdPolicy, error) {
	policies, err := listNetbirdPolicies(ctx, ListNetbirdPoliciesParams{})
	if err != nil {
		return nil, err
	}

	group := strings.ToLower(strings.TrimSpace(args.Group))
	if group == "" {
		return nil, fmt.Errorf("group must not be empty")
	}

	var matches []NetbirdPolicy
	for _, policy := range policies {
		var rules []NetbirdPolicyRule
		for _, rule := range policy.Rules {
			for _, source := range rule.Sources {
				if strings.EqualFold(source.Name, group) || strings.EqualFold(source.ID, group) {
					rules = append(rules, rule)
					break
				}
			}
			for _, destination := range rule.Destinations {
				if strings.EqualFold(destination.Name, group) || strings.EqualFold(destination.ID, group) {
					rules = append(rules, rule)
					break
				}
			}
		}
		if len(rules) > 0 {
			policy.Rules = rules
			matches = append(matches, policy)
		}
	}

	return matches, nil
}

var FindNetbirdPolicies = mcpnetbird.MustTool(
	"find_netbird_policies",
	"Find Netbird policies with a source or destination group name or ID. Returns only matching rules.",
	findNetbirdPolicies,
)

func AddNetbirdPolicyTools(mcp *server.MCPServer) {
	ListNetbirdPolicies.Register(mcp)
	FindNetbirdPolicies.Register(mcp)
}
