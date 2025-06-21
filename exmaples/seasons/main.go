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

	// List all seasons
	seasons, err := client.Seasons.List(ctx)
	if err != nil {
		log.Fatalln("Failed to retrieve seasons list:", err)
	}

	fmt.Println("=== Seasons ===")
	for _, season := range seasons.Data {
		fmt.Printf("- %s (Id: %d)\n", season.Name, season.Id)
	}

	// Get the current season's active machine
	activeMachine, err := client.Seasons.ActiveMachine(ctx)
	if err != nil {
		log.Fatalln("Failed to get active machine:", err)
	}

	fmt.Println("\n=== Active Machine ===")
	fmt.Printf("%s (Id: %d)\n", activeMachine.Data.Name, activeMachine.Data.Id)

	// Optional: Print raw response
	fmt.Println("\nRaw Response:", string(activeMachine.ResponseMeta.Raw))
}
