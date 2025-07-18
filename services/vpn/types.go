package vpn

import (
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

type ConnectionStatusItem struct {
	Connection           Connection
	LocationTypeFriendly string
	Server               ConnectionStatusServer
	Type                 string
}

type ServerConnection struct {
	CurrentClients int
	FriendlyName   string
	Id             int
	Location       string
}

type Connection struct {
	Down          float32
	Ip4           string
	Ip6           string
	Name          string
	ThroughPwnbox bool
	Up            float32
}

type ConnectionStatusServer struct {
	FriendlyName string
	Hostname     string
	Id           int
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

type AssignedServerConnectionsServers struct {
	CurrentClients       int
	FriendlyName         string
	Id                   int
	Location             string
	LocationTypeFriendly string
}

type Server struct {
	CurrentClients int
	FriendlyName   string
	Full           bool
	Id             int
	Location       string
	Tier           string
}
