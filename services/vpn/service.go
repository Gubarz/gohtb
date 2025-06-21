package vpn

import (
	"context"
	"sort"
	"strings"

	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/convert"
	"github.com/gubarz/gohtb/internal/deref"
	"github.com/gubarz/gohtb/internal/errutil"
	"github.com/gubarz/gohtb/internal/extract"
	v4client "github.com/gubarz/gohtb/internal/httpclient/v4"
	"github.com/gubarz/gohtb/internal/ptr"
	"github.com/gubarz/gohtb/internal/service"
)

func NewService(client service.Client) *Service {
	return &Service{
		base: service.NewBase(client),
	}
}

// VPN returns a handle for a specific VPN endpoints
// with the given ID. This handle can be used to perform
// operations related to that VPN endpoint, such as
// downloading configuration files or switching servers.
// The ID is typically obtained from the VPN Status
// or from the Server.Data.Options results.
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
//	ctx := context.Background()
//	resp, err := htb.VPN.VPN(256).DownloadUDP(ctx)
//	fmt.Println("VPN file:", string(resp.Data))
func (h *Handle) DownloadUDP(ctx context.Context) (VPNFileResponse, error) {
	resp, err := h.client.V4().GetAccessOvpnfileVpnIdUDPWithResponse(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil {
		return errutil.UnwrapFailure(err, raw, resp.StatusCode(), func(raw []byte) VPNFileResponse {
			return VPNFileResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return VPNFileResponse{
		Data: raw,
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
			CFRay:      resp.HTTPResponse.Header.Get("CF-Ray"),
		},
	}, nil
}

func (h *Handle) DownloadTCP(ctx context.Context) (VPNFileResponse, error) {
	resp, err := h.client.V4().GetAccessOvpnfileVpnIdTCPWithResponse(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil {
		return errutil.UnwrapFailure(err, raw, resp.StatusCode(), func(raw []byte) VPNFileResponse {
			return VPNFileResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return VPNFileResponse{
		Data: raw,
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
			CFRay:      resp.HTTPResponse.Header.Get("CF-Ray"),
		},
	}, nil
}

func (s *Service) Status(ctx context.Context) (ConnectionStatusResponse, error) {
	resp, err := s.base.Client.V4().GetConnectionStatusWithResponse(s.base.Client.Limiter().Wrap(ctx))

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, resp.StatusCode(), func(raw []byte) ConnectionStatusResponse {
			return ConnectionStatusResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return ConnectionStatusResponse{
		Data: convert.Slice(*resp.JSON200, fromAPIConnectionStatusItem),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
			CFRay:      resp.HTTPResponse.Header.Get("CF-Ray"),
		},
	}, nil
}

func (s *Service) Connections(ctx context.Context) (ConnectionStatusResponse, error) {
	resp, err := s.base.Client.V4().GetConnectionStatusWithResponse(s.base.Client.Limiter().Wrap(ctx))

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, resp.StatusCode(), func(raw []byte) ConnectionStatusResponse {
			return ConnectionStatusResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return ConnectionStatusResponse{
		Data: convert.Slice(*resp.JSON200, fromAPIConnectionStatusItem),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
			CFRay:      resp.HTTPResponse.Header.Get("CF-Ray"),
		},
	}, nil
}

func (s *Service) Servers(product string) *ServerQuery {
	return &ServerQuery{
		client:  s.base.Client,
		product: product,
	}
}

func (q *ServerQuery) ByTier(tier string) *ServerQuery {
	qc := ptr.Clone(q)
	qc.tier = tier
	return qc
}

func (q *ServerQuery) ByLocation(location string) *ServerQuery {
	qc := ptr.Clone(q)
	qc.location = location
	return qc
}

func (q *ServerQuery) Results(ctx context.Context) (ConnectionsServersResponse, error) {
	resp, err := q.client.V4().GetConnectionsServersWithResponse(q.client.Limiter().Wrap(ctx),
		&v4client.GetConnectionsServersParams{
			Product: v4client.GetConnectionsServersParamsProduct(q.product),
		})

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, resp.StatusCode(), func(raw []byte) ConnectionsServersResponse {
			return ConnectionsServersResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	var flat []Server
	if resp.JSON200.Data != nil && resp.JSON200.Data.Options != nil {
		flat = flattenOptions(resp.JSON200.Data.Options, q.tier, q.location)
	}

	var assigned AssignedServerConnectionsServers
	if resp.JSON200.Data.Assigned != nil {
		assigned = fromAPIAssignedServerConnectionsServers(resp.JSON200.Data.Assigned)
	}

	res := ConnectionsServerData{
		Assigned: assigned,
		Disabled: deref.Bool(resp.JSON200.Data.Disabled),
		Options:  flat,
	}

	return ConnectionsServersResponse{
		Data: res,
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
			CFRay:      resp.HTTPResponse.Header.Get("CF-Ray"),
		},
	}, nil
}

func fromAPIAssignedServerConnectionsServers(data *v4client.AssignedServerConnectionsServers) AssignedServerConnectionsServers {
	return AssignedServerConnectionsServers{
		CurrentClients:       deref.Int(data.CurrentClients),
		FriendlyName:         deref.String(data.FriendlyName),
		Id:                   deref.Int(data.Id),
		Location:             deref.String(data.Location),
		LocationTypeFriendly: deref.String(data.LocationTypeFriendly),
	}
}

func flattenOptions(opts *v4client.Options, tier, location string) []Server {
	tier = strings.ToLower(strings.TrimSpace(tier))
	location = strings.ToLower(strings.TrimSpace(location))

	var out []Server

	for _, regionMap := range *opts {
		for _, group := range regionMap {
			if group.Servers == nil {
				continue
			}
			for _, server := range *group.Servers {
				loc := deref.String(server.Location)
				name := deref.String(server.FriendlyName)

				actualTier := extractTierFromFriendly(name)

				if location != "" && !matchLocationField(loc, location) {
					continue
				}
				if tier != "" && actualTier != tier {
					continue
				}

				out = append(out, Server{
					Id:             deref.Int(server.Id),
					FriendlyName:   name,
					Location:       loc,
					CurrentClients: deref.Int(server.CurrentClients),
					Full:           deref.Bool(server.Full),
					Tier:           actualTier,
				})
			}
		}
	}

	return out
}

func matchLocationField(locationField, filter string) bool {
	return strings.EqualFold(strings.TrimSpace(locationField), strings.TrimSpace(filter))
}

func extractTierFromFriendly(name string) string {
	name = strings.ToLower(name)
	switch {
	case strings.Contains(" "+name+" ", " vip+ "):
		return "vip+"
	case strings.Contains(" "+name+" ", " vip "):
		return "vip"
	case strings.Contains(" "+name+" ", " free "):
		return "free"
	default:
		return "unknown"
	}
}

func (h *Handle) Switch(ctx context.Context) (common.MessageResponse, error) {
	resp, err := h.client.V4().PostConnectionsServersSwitchWithResponse(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil {
		return errutil.UnwrapFailure(err, raw, resp.StatusCode(), func(raw []byte) common.MessageResponse {
			return common.MessageResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return common.MessageResponse{
		Data: common.Message{
			Message: deref.String(resp.JSON200.Message),
		},
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
			CFRay:      resp.HTTPResponse.Header.Get("CF-Ray"),
		},
	}, nil
}

func (o OptionsServers) ByLocation(location string) OptionsServers {
	var d OptionsServers
	for _, v := range o {
		if strings.EqualFold(v.Location, location) {
			d = append(d, v)
		}
	}
	return d
}

func (o OptionsServers) ByTier(tier string) OptionsServers {
	var d OptionsServers
	for _, v := range o {
		if strings.EqualFold(v.Tier, tier) {
			d = append(d, v)
		}
	}
	return d
}

func (o OptionsServers) SortByCurrentClients() OptionsServers {
	sort.Slice(o, func(i, j int) bool {
		return o[i].CurrentClients < o[j].CurrentClients
	})
	return o
}

func (o OptionsServers) First() Server {
	if len(o) == 0 {
		return Server{}
	}
	return o[0]
}
