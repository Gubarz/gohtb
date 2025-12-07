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

	// List all fortresses
	fortresses, err := client.Fortresses.List(ctx)
	if err != nil {
		log.Fatalln("Failed to retrieve fortress list:", err)
	}

	fmt.Println("=== Fortresses ===")
	for _, fortress := range fortresses.Data {
		fmt.Printf("- %s (Id: %d)\n", fortress.Name, fortress.Id)
	}

	if len(fortresses.Data) > 0 {
		// Get details for the first fortress
		firstFortress := fortresses.Data[0]
		fortressDetails, err := client.Fortresses.Fortress(firstFortress.Id).Info(ctx)
		if err != nil {
			log.Printf("Failed to get fortress details for %s: %v\n", firstFortress.Name, err)
		} else {
			fmt.Printf("\n=== Fortress Details: %s ===\n", fortressDetails.Data.Name)
			fmt.Printf("Description: %s\n", fortressDetails.Data.Description)
			fmt.Printf("Points: %s\n", fortressDetails.Data.Points)

			// Optional: Print raw response
			fmt.Println("\nRaw Response:", string(fortressDetails.ResponseMeta.Raw))
		}
	} else {
		fmt.Println("\nNo fortresses found.")
	}
}
