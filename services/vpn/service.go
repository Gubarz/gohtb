package vpn

import (
	"context"
	"sort"
	"strconv"
	"strings"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	v5Client "github.com/gubarz/gohtb/httpclient/v5"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/errutil"
	"github.com/gubarz/gohtb/internal/extract"
	"github.com/gubarz/gohtb/internal/ptr"
	"github.com/gubarz/gohtb/internal/service"
)

const (
	Competitive   = v4Client.ProductCompetitive
	Fortresses    = v4Client.ProductFortresses
	Labs          = v4Client.ProductLabs
	StartingPoint = v4Client.ProductStartingPoint
)

type ServerQuery struct {
	client   service.Client
	product  string
	tier     string
	location string
}

type ProlabQuery struct {
	client   service.Client
	prolab   int
	tier     string
	location string
}

type VPNFileResponse struct {
	Data         []byte
	ResponseMeta common.ResponseMeta
}

type ConnectionStatusItem = v4Client.ConnectionStatusItem

type ConnectionStatusResponse struct {
	Data         []ConnectionStatusItem
	ResponseMeta common.ResponseMeta
}

type ConnectionsServersResponse struct {
	Data         ConnectionsServerData
	ResponseMeta common.ResponseMeta
}

type AssignedServerConnectionsServers = v4Client.AssignedServerConnectionsServers

type OptionsServers []Server

type ConnectionsServerData struct {
	Assigned AssignedServerConnectionsServers
	Disabled bool

	Options OptionsServers
}

type Server struct {
	CurrentClients int
	FriendlyName   string
	Full           bool
	Id             int
	Location       string
	VpnType        string
	VpnNumber      int
}

type Service struct {
	base service.Base
}

// NewService creates a new VPN service bound to a shared client.
//
// Example:
//
//	vpnService := vpn.NewService(client)
//	_ = vpnService
func NewService(client service.Client) *Service {
	return &Service{
		base: service.NewBase(client),
	}
}

type Handle struct {
	client service.Client
	id     int
}

// VPN returns a handle for a specific VPN endpoints
// with the given ID. This handle can be used to perform
// operations related to that VPN endpoint, such as
// downloading configuration files or switching servers.
// The ID is typically obtained from the VPN Status
// or from the Server.Data.Options results.
//
// Example:
//
//	endpoint := client.VPN.VPN(256)
//	_ = endpoint
func (s *Service) VPN(id int) *Handle {
	return &Handle{
		client: s.base.Client,
		id:     id,
	}
}

// DownloadUDP downloads the OpenVPN configuration file for the specified VPN endpoint using UDP.
//
// Example:
//
//	resp, err := client.VPN.VPN(256).DownloadUDP(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Downloaded UDP config: %d bytes\n", len(resp.Data))
func (h *Handle) DownloadUDP(ctx context.Context) (VPNFileResponse, error) {
	resp, err := h.client.V4().GetAccessOvpnfileVpnIdUDP(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || len(raw) == 0 {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) VPNFileResponse {
			return VPNFileResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return VPNFileResponse{
		Data: raw,
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode,
			Headers:    resp.Header,
			CFRay:      resp.Header.Get("CF-Ray"),
		},
	}, nil
}

// DownloadTCP downloads the OpenVPN configuration file for the specified VPN endpoint using TCP.
//
// Example:
//
//	resp, err := client.VPN.VPN(256).DownloadTCP(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Downloaded TCP config: %d bytes\n", len(resp.Data))
func (h *Handle) DownloadTCP(ctx context.Context) (VPNFileResponse, error) {
	resp, err := h.client.V4().GetAccessOvpnfileVpnIdTCP(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || len(raw) == 0 {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) VPNFileResponse {
			return VPNFileResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return VPNFileResponse{
		Data: raw,
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode,
			Headers:    resp.Header,
			CFRay:      resp.Header.Get("CF-Ray"),
		},
	}, nil
}

// Status retrieves the current VPN connection status information.
// This includes details about active connections and their current state.
//
// Example:
//
//	status, err := client.VPN.Status(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Connection status: %+v\n", status.Data)
func (s *Service) Status(ctx context.Context) (ConnectionStatusResponse, error) {
	resp, err := s.base.Client.V4().GetConnectionStatus(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return ConnectionStatusResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetConnectionStatusResponse)
	if err != nil {
		return ConnectionStatusResponse{ResponseMeta: meta}, err
	}

	return ConnectionStatusResponse{
		Data:         *parsed.JSON200,
		ResponseMeta: meta,
	}, nil
}

// Servers creates a new server query for the specified product.
// This returns a ServerQuery that can be chained with filtering methods
// like ByTier() and ByLocation() before calling Results().
//
// Common products include "labs", "starting_point", "fortresses", etc.
//
// Example:
//
//	query := client.VPN.Servers("labs")
//	servers, err := query.ByTier("free").ByLocation("US").Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Matching servers: %d\n", len(servers.Data.Options))
func (s *Service) Servers(product string) *ServerQuery {
	return &ServerQuery{
		client:  s.base.Client,
		product: product,
	}
}

// ByTier filters the server query by tier using case-insensitive matching.
// Common values include "free", "dedivip", "release", "fort", and "starting-point".
// Returns a new ServerQuery that can be further chained.
//
// Example:
//
//	servers, err := client.VPN.Servers("labs").ByTier("free").Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Free servers: %d\n", len(servers.Data.Options))
func (q *ServerQuery) ByTier(tier string) *ServerQuery {
	qc := ptr.Clone(q)
	qc.tier = tier
	return qc
}

// ByLocation filters the server query by location using case-insensitive matching.
// Returns a new ServerQuery that can be further chained.
//
// Example:
//
//	servers, err := client.VPN.Servers("labs").ByLocation("US").Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("US servers: %d\n", len(servers.Data.Options))
func (q *ServerQuery) ByLocation(location string) *ServerQuery {
	qc := ptr.Clone(q)
	qc.location = location
	return qc
}

// Results executes the server query and returns the filtered server list.
// This method should be called last in the query chain to fetch the actual data.
// The returned servers are flattened from the API's nested structure and include
// tier information extracted from server names.
//
// Example:
//
//	servers, err := client.VPN.Servers("labs").
//		ByTier("free").
//		ByLocation("US").
//		Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, server := range servers.Data.Options {
//		fmt.Printf("Server: %s (%s)\n", server.FriendlyName, server.Location)
//	}
func (q *ServerQuery) Results(ctx context.Context) (ConnectionsServersResponse, error) {
	resp, err := q.client.V4().GetConnectionsServers(q.client.Limiter().Wrap(ctx),
		&v4Client.GetConnectionsServersParams{
			Product: v4Client.GetConnectionsServersParamsProduct(q.product),
		})

	if err != nil {
		return ConnectionsServersResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetConnectionsServersResponse)
	if err != nil {
		return ConnectionsServersResponse{ResponseMeta: meta}, err
	}

	var flat []Server
	var assigned AssignedServerConnectionsServers

	if parsed != nil && parsed.JSON200 != nil {
		data := parsed.JSON200.Data

		if data.Options != nil {
			flat = flattenOptions(data.Options, q.tier, q.location)
		}

		assigned = data.Assigned
	}

	res := ConnectionsServerData{
		Assigned: assigned,
		Disabled: parsed.JSON200.Data.Disabled,
		Options:  flat,
	}

	return ConnectionsServersResponse{
		Data:         res,
		ResponseMeta: meta,
	}, nil
}

func flattenOptions(opts v4Client.Options, tier, location string) []Server {
	tier = strings.ToLower(strings.TrimSpace(tier))
	location = strings.ToLower(strings.TrimSpace(location))

	var out []Server

	for _, regionMap := range opts {
		for _, group := range regionMap {
			if group.Servers == nil {
				continue
			}
			for _, server := range group.Servers {
				loc := server.Location
				name := server.FriendlyName

				vpnType, vpnNumber := extractTierFromFriendly(name)

				if location != "" && !matchLocationField(loc, location) {
					continue
				}
				if tier != "" && vpnType != tier {
					continue
				}

				out = append(out, Server{
					Id:             server.Id,
					FriendlyName:   name,
					Location:       loc,
					CurrentClients: server.CurrentClients,
					Full:           server.Full,
					VpnType:        vpnType,
					VpnNumber:      vpnNumber,
				})
			}
		}
	}

	return out
}

func matchLocationField(locationField, filter string) bool {
	return strings.EqualFold(strings.TrimSpace(locationField), strings.TrimSpace(filter))
}

func extractTierFromFriendly(name string) (string, int) {
	var vpnType string
	var vpnNumber int
	name = strings.ToLower(name)
	switch {
	case strings.Contains(name, "vip+"):
		vpnType = "dedivip"
	case strings.Contains(name, "release arena"):
		vpnType = "release"
	case strings.Contains(name, "fortress"):
		vpnType = "fort"
	case strings.Contains(name, "startingpoint"):
		vpnType = "starting-point"
	case strings.Contains(name, "mini pro lab free"):
		vpnType = "mini-prolab-free"
	default:
		// Servers without specific markers are free tier
		vpnType = "free"
	}
	// Attempt to extract last part which is the number
	parts := strings.Split(name, " ")
	lastPart := parts[len(parts)-1]
	vpnNumber, err := strconv.Atoi(lastPart)
	if err != nil {
		vpnNumber = 0
	}

	return vpnType, vpnNumber
}

// Switch changes the VPN connection to the server specified by this handle's ID.
// Returns a message response indicating the result of the switch operation.
//
// Example:
//
//	result, err := client.VPN.VPN(256).Switch(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Switch result:", result.Data.Message)
func (h *Handle) Switch(ctx context.Context) (common.MessageResponse, error) {
	resp, err := h.client.V4().PostConnectionsServersSwitch(h.client.Limiter().Wrap(ctx), h.id)
	if err != nil {
		return common.MessageResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParsePostVMTerminateResponse)
	if err != nil {
		return common.MessageResponse{ResponseMeta: meta}, err
	}

	return common.MessageResponse{
		Data: common.Message{
			Message: parsed.JSON200.Message,
		},
		ResponseMeta: meta,
	}, nil
}

// ByLocation filters servers by location using case-insensitive matching.
// Returns a new OptionsServers slice containing only servers in the specified location.
//
// Example:
//
//	servers := response.Data.Options.ByLocation("US")
func (o OptionsServers) ByLocation(location string) OptionsServers {
	var d OptionsServers
	for _, v := range o {
		if strings.EqualFold(v.Location, location) {
			d = append(d, v)
		}
	}
	return d
}

// ByType filters servers by VPN type using case-insensitive matching.
// Valid types include "free", "starting-point", "fort", "dedivip", and "release".
// Returns a new OptionsServers slice containing only servers of the specified type.
//
// Example:
//
//	freeServers := response.Data.Options.ByType("free")
//	dedicatedVIPServers := response.Data.Options.ByType("dedivip")
func (o OptionsServers) ByType(vpnType string) OptionsServers {
	var d OptionsServers
	for _, v := range o {
		if strings.EqualFold(v.VpnType, vpnType) {
			d = append(d, v)
		}
	}
	return d
}

// SortByCurrentClients sorts servers by current client count in ascending order.
// Servers with fewer current clients appear first in the returned slice.
// This method modifies the original slice and returns it for method chaining.
//
// Example:
//
//	leastBusyFirst := response.Data.Options.SortByCurrentClients()
func (o OptionsServers) SortByCurrentClients() OptionsServers {
	sort.Slice(o, func(i, j int) bool {
		return o[i].CurrentClients < o[j].CurrentClients
	})
	return o
}

// First returns the first server in the slice, or an empty Server if the slice is empty.
// This is commonly used after filtering and sorting to get the best match.
//
// Example:
//
//	bestServer := response.Data.Options.
//		ByType("free").
//		ByLocation("US").
//		SortByCurrentClients().
//		First()
func (o OptionsServers) First() Server {
	if len(o) == 0 {
		return Server{}
	}
	return o[0]
}

func (h *Handle) switchAndDownload(ctx context.Context, useUDP bool) (VPNFileResponse, error) {
	resp, err := h.Switch(ctx)

	if err != nil {
		return errutil.UnwrapFailure(err, resp.Raw, common.SafeStatus(resp), func(raw []byte) VPNFileResponse {
			return VPNFileResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	if useUDP {
		return h.DownloadUDP(ctx)
	}
	return h.DownloadTCP(ctx)
}

// SwitchAndDownloadUDP is a convenience method that switches servers and downloads UDP config.
//
// Example:
//
//	resp, err := client.VPN.VPN(256).SwitchAndDownloadUDP(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Downloaded UDP config: %d bytes\n", len(resp.Data))
func (h *Handle) SwitchAndDownloadUDP(ctx context.Context) (VPNFileResponse, error) {
	return h.switchAndDownload(ctx, true)
}

// SwitchAndDownloadTCP is a convenience method that switches servers and downloads TCP config.
//
// Example:
//
//	resp, err := client.VPN.VPN(256).SwitchAndDownloadTCP(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Downloaded TCP config: %d bytes\n", len(resp.Data))
func (h *Handle) SwitchAndDownloadTCP(ctx context.Context) (VPNFileResponse, error) {
	return h.switchAndDownload(ctx, false)
}

// ProlabServers creates a prolab VPN server query for a specific prolab ID.
//
// Example:
//
//	query := client.VPN.ProlabServers(8)
//	_ = query
func (s *Service) ProlabServers(id int) *ProlabQuery {
	return &ProlabQuery{
		client: s.base.Client,
		prolab: id,
	}
}

// ByTier filters the server query by tier using case-insensitive matching.
// Common values include "free", "dedivip", "release", "fort", and "starting-point".
// Returns a new ServerQuery that can be further chained.
//
// Example:
//
//	servers, err := client.VPN.ProlabServers(8).ByTier("free").Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Prolab servers: %d\n", len(servers.Data.Options))
func (q *ProlabQuery) ByTier(tier string) *ProlabQuery {
	qc := ptr.Clone(q)
	qc.tier = tier
	return qc
}

// ByLocation filters the server query by location using case-insensitive matching.
// Returns a new ServerQuery that can be further chained.
//
// Example:
//
//	servers, err := client.VPN.ProlabServers(8).ByLocation("US").Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("US prolab servers: %d\n", len(servers.Data.Options))
func (q *ProlabQuery) ByLocation(location string) *ProlabQuery {
	qc := ptr.Clone(q)
	qc.location = location
	return qc
}

// Results executes the server query and returns the filtered server list.
// This method should be called last in the query chain to fetch the actual data.
// The returned servers are flattened from the API's nested structure and include
// tier information extracted from server names.
//
// Example:
//
//	servers, err := client.VPN.ProlabServers(8).
//		ByTier("free").
//		ByLocation("US").
//		Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, server := range servers.Data.Options {
//		fmt.Printf("Server: %s (%s)\n", server.FriendlyName, server.Location)
//	}
func (q *ProlabQuery) Results(ctx context.Context) (ConnectionsServersResponse, error) {
	resp, err := q.client.V4().GetConnectionsServersProlab(q.client.Limiter().Wrap(ctx),
		q.prolab)

	if err != nil {
		return ConnectionsServersResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetConnectionsServersResponse)
	if err != nil {
		return ConnectionsServersResponse{ResponseMeta: meta}, err
	}

	var flat []Server
	var assigned AssignedServerConnectionsServers

	if parsed != nil && parsed.JSON200 != nil {
		data := parsed.JSON200.Data

		if data.Options != nil {
			flat = flattenOptions(data.Options, q.tier, q.location)
		}

		assigned = data.Assigned
	}

	res := ConnectionsServerData{
		Assigned: assigned,
		Disabled: parsed.JSON200.Data.Disabled,
		Options:  flat,
	}

	return ConnectionsServersResponse{
		Data:         res,
		ResponseMeta: meta,
	}, nil
}

// ProductName identifies product names accepted by the connection status endpoint.
type ProductName = v4Client.GetConnectionStatusProductnameParamsProductName

const (
	ProductNameCompetitive   ProductName = v4Client.GetConnectionStatusProductnameParamsProductNameCompetitive
	ProductNameFortresses    ProductName = v4Client.GetConnectionStatusProductnameParamsProductNameFortresses
	ProductNameLabs          ProductName = v4Client.GetConnectionStatusProductnameParamsProductNameLabs
	ProductNameStartingPoint ProductName = v4Client.GetConnectionStatusProductnameParamsProductNameStartingPoint
)

type ProductConnectionStatusData = v4Client.ConnectionStatusProductResponse

// ProductConnectionStatusResponse contains connection status data for a specific product/prolab.
type ProductConnectionStatusResponse struct {
	Data         ProductConnectionStatusData
	ResponseMeta common.ResponseMeta
}

// StatusByProlab retrieves connection status for a specific prolab.
//
// Example:
//
//	status, err := client.VPN.StatusByProlab(ctx, 1)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Prolab connection status: %+v\n", status.Data)
func (s *Service) StatusByProlab(ctx context.Context, prolabId int) (ProductConnectionStatusResponse, error) {
	resp, err := s.base.Client.V4().GetConnectionStatusProlab(s.base.Client.Limiter().Wrap(ctx), prolabId)
	if err != nil {
		return ProductConnectionStatusResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetConnectionStatusProlabResponse)
	if err != nil {
		return ProductConnectionStatusResponse{ResponseMeta: meta}, err
	}

	return ProductConnectionStatusResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

// StatusByProductName retrieves connection status for a specific product name.
//
// Example:
//
//	status, err := client.VPN.StatusByProductName(ctx, vpn.ProductNameLabs)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Product connection status: %+v\n", status.Data)
func (s *Service) StatusByProductName(ctx context.Context, productName ProductName) (ProductConnectionStatusResponse, error) {
	resp, err := s.base.Client.V4().GetConnectionStatusProductname(
		s.base.Client.Limiter().Wrap(ctx),
		v4Client.GetConnectionStatusProductnameParamsProductName(productName),
	)
	if err != nil {
		return ProductConnectionStatusResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetConnectionStatusProductnameResponse)
	if err != nil {
		return ProductConnectionStatusResponse{ResponseMeta: meta}, err
	}

	return ProductConnectionStatusResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

type ConnectionsV4Data = v4Client.ConnectionsResponse

// ConnectionsV4Response contains v4 connections payload.
type ConnectionsV4Response struct {
	Data         ConnectionsV4Data
	ResponseMeta common.ResponseMeta
}

// ConnectionsDeprecated retrieves v4 connection inventory from /connections.
//
// Example:
//
//	connections, err := client.VPN.ConnectionsDeprecated(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("V4 connections payload: %+v\n", connections.Data)
func (s *Service) ConnectionsDeprecated(ctx context.Context) (ConnectionsV4Response, error) {
	resp, err := s.base.Client.V4().GetConnections(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return ConnectionsV4Response{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetConnectionsResponse)
	if err != nil {
		return ConnectionsV4Response{ResponseMeta: meta}, err
	}

	return ConnectionsV4Response{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

type ConnectionsV5Data = v5Client.ConnectionsResponse

// ConnectionsV5Response contains v5 connections payload.
type ConnectionsV5Response struct {
	Data         ConnectionsV5Data
	ResponseMeta common.ResponseMeta
}

// Connections retrieves v5 connection inventory from /connections.
//
// Example:
//
//	connections, err := client.VPN.Connections(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("V5 connections payload: %+v\n", connections.Data)
func (s *Service) Connections(ctx context.Context) (ConnectionsV5Response, error) {
	resp, err := s.base.Client.V5().GetConnections(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return ConnectionsV5Response{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v5Client.ParseGetConnectionsResponse)
	if err != nil {
		return ConnectionsV5Response{ResponseMeta: meta}, err
	}

	return ConnectionsV5Response{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}
