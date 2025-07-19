package vpn

import (
	v4client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/deref"
)

func fromAPIConnectionStatusItem(data v4client.ConnectionStatusItem) ConnectionStatusItem {
	return ConnectionStatusItem{
		Connection:           fromAPIConnection(data.Connection),
		LocationTypeFriendly: deref.String(data.LocationTypeFriendly),
		Server:               fromAPIConnectionStatusServer(data.Server),
		Type:                 deref.String(data.Type),
	}
}

func fromAPIConnection(data *v4client.Connection) Connection {
	if data == nil {
		return Connection{}
	}
	return Connection{
		Down:          deref.Float32(data.Down),
		Ip4:           deref.String(data.Ip4),
		Ip6:           deref.String(data.Ip6),
		Name:          deref.String(data.Name),
		ThroughPwnbox: deref.Bool(data.ThroughPwnbox),
		Up:            deref.Float32(data.Up),
	}
}

func fromAPIConnectionStatusServer(data *v4client.ConnectionServer) ConnectionStatusServer {
	if data == nil {
		return ConnectionStatusServer{}
	}
	return ConnectionStatusServer{
		FriendlyName: deref.String(data.FriendlyName),
		Hostname:     deref.String(data.Hostname),
		Id:           deref.Int(data.Id),
		ProLabId:     deref.Int(data.ProLabId),
	}
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
