package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gubarz/gohtb"
)

func main() {
	ctx := context.Background()

	client, err := gohtb.New(os.Getenv("HTB_TOKEN"))
	if err != nil {
		log.Fatalln("Failed to create HTB client:", err)
	}

	// Find the best free US server
	fmt.Println("\n=== Finding Best Free US Server ===")
	servers, err := client.VPN.Servers("labs").ByTier("free").ByLocation("US").Results(ctx)
	if err != nil {
		log.Printf("Failed to get servers: %v\n", err)
	} else {
		bestServer := servers.Data.Options.SortByCurrentClients().First()
		if bestServer.Id != 0 {
			fmt.Printf("Best server: %s (%d clients)\n", bestServer.FriendlyName, bestServer.CurrentClients)

			// Switch to the best server and download UDP config
			fmt.Println("\n=== Switching Server and Downloading Config ===")
			vpnConfig, err := client.VPN.VPN(bestServer.Id).SwitchAndDownloadUDP(ctx)
			if err != nil {
				log.Printf("Failed to switch and download: %v\n", err)
			} else {
				filename := fmt.Sprintf("%s.ovpn", bestServer.FriendlyName)
				err = os.WriteFile(filename, vpnConfig.Data, 0644)
				if err != nil {
					log.Printf("Failed to write config file: %v\n", err)
				} else {
					fmt.Printf("VPN config saved to: %s (%d bytes)\n", filename, len(vpnConfig.Data))
				}
			}
		} else {
			fmt.Println("No free US servers found")
		}
	}
}
