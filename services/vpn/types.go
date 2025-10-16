package vpn

import (
	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/service"
)

type Service struct {
	base service.Base
}

type Handle struct {
	client service.Client
	id     int
}

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
type ConnectionStatusResponse struct {
	Data         []ConnectionStatusItem
	ResponseMeta common.ResponseMeta
}

type ServerConnection struct {
	CurrentClients int
	FriendlyName   string
	Id             int
	Location       string
}

type ConnectionsServersResponse struct {
	Data         ConnectionsServerData
	ResponseMeta common.ResponseMeta
}

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
	Tier           string
}

type ConnectionStatusItem = v4Client.ConnectionStatusItem
type Connection = v4Client.Connection
type ConnectionStatusServer = v4Client.ConnectionServer
type AssignedServerConnectionsServers = v4Client.AssignedServerConnectionsServers
