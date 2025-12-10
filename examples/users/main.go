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

	// Example user ID - replace with actual user ID
	userID := 1337

	// Get user profile information
	fmt.Printf("=== User Profile (ID: %d) ===\n", userID)
	profile, err := client.Users.User(userID).ProfileBasic(ctx)
	if err != nil {
		log.Printf("Failed to get user profile: %v\n", err)
	} else {
		fmt.Printf("Username: %s\n", profile.Data.Name)
		fmt.Printf("Points: %d\n", profile.Data.Points)
		fmt.Printf("Country: %s\n", profile.Data.CountryName)
		fmt.Printf("Owns: User: %d, Root: %d\n", profile.Data.UserOwns, profile.Data.SystemOwns)
	}

	// Get user activity
	fmt.Printf("\n=== User Activity ===\n")
	activity, err := client.Users.User(userID).ProfileActivity().Results(ctx)
	if err != nil {
		log.Printf("Failed to get user activity: %v\n", err)
	} else {
		fmt.Printf("Recent activities: %d entries\n", len(activity.Data))
		for i, act := range activity.Data {
			if i >= 5 { // Show only first 5 activities
				break
			}
			fmt.Printf("- %s: %s (%s)\n", act.OwnDate, act.Type, act.Name)
		}
	}
}
