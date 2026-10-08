package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	mcpnetbird "github.com/aantti/mcp-netbird"
	"github.com/mark3labs/mcp-go/server"
)

type NetbirdPeerGroup struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	PeersCount     int    `json:"peers_count"`
	ResourcesCount int    `json:"resources_count"`
}

type NetbirdPeer struct {
	AccessiblePeersCount        int                `json:"accessible_peers_count"`
	ApprovalRequired            bool               `json:"approval_required"`
	CityName                    string             `json:"city_name"`
	Connected                   bool               `json:"connected"`
	ConnectionIP                string             `json:"connection_ip"`
	CountryCode                 string             `json:"country_code"`
	DNSLabel                    string             `json:"dns_label"`
	ExtraDNSLabels              []string           `json:"extra_dns_labels"`
	GeonameID                   int                `json:"geoname_id"`
	Groups                      []NetbirdPeerGroup `json:"groups"`
	Hostname                    string             `json:"hostname"`
	ID                          string             `json:"id"`
	InactivityExpirationEnabled bool               `json:"inactivity_expiration_enabled"`
	IP                          string             `json:"ip"`
	KernelVersion               string             `json:"kernel_version"`
	LastLogin                   time.Time          `json:"last_login"`
	LastSeen                    time.Time          `json:"last_seen"`
	LoginExpirationEnabled      bool               `json:"login_expiration_enabled"`
	LoginExpired                bool               `json:"login_expired"`
	Name                        string             `json:"name"`
	OS                          string             `json:"os"`
	SerialNumber                string             `json:"serial_number"`
	SSHEnabled                  bool               `json:"ssh_enabled"`
	UIVersion                   string             `json:"ui_version"`
	UserID                      string             `json:"user_id"`
	Version                     string             `json:"version"`
}

type ListNetbirdPeersParams struct{}

func listNetbirdPeers(ctx context.Context, args ListNetbirdPeersParams) ([]NetbirdPeer, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient()
	}

	var peers []NetbirdPeer
	if err := client.Get(ctx, "/peers", &peers); err != nil {
		return nil, err
	}

	return peers, nil
}

var ListNetbirdPeers = mcpnetbird.MustTool(
	"list_netbird_peers",
	"List all Netbird peers",
	listNetbirdPeers,
)

type FindNetbirdPeersParams struct {
	Query string `mcp:"query" validate:"required" json:"query"`
}

func findNetbirdPeers(ctx context.Context, args FindNetbirdPeersParams) ([]NetbirdPeer, error) {
	peers, err := listNetbirdPeers(ctx, ListNetbirdPeersParams{})
	if err != nil {
		return nil, err
	}

	query := strings.ToLower(strings.TrimSpace(args.Query))
	if query == "" {
		return nil, fmt.Errorf("query must not be empty")
	}

	var matches []NetbirdPeer
	for _, peer := range peers {
		if strings.Contains(strings.ToLower(peer.Name), query) ||
			strings.Contains(strings.ToLower(peer.Hostname), query) ||
			strings.Contains(strings.ToLower(peer.DNSLabel), query) {
			matches = append(matches, peer)
		}
	}

	return matches, nil
}

var FindNetbirdPeers = mcpnetbird.MustTool(
	"find_netbird_peers",
	"Find Netbird peers by name, hostname, or DNS label. Use this instead of listing all peers when diagnosing one user.",
	findNetbirdPeers,
)

func AddNetbirdPeerTools(mcp *server.MCPServer) {
	ListNetbirdPeers.Register(mcp)
	FindNetbirdPeers.Register(mcp)
}
