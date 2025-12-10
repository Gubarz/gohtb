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

	// Example team ID - replace with actual team ID
	teamID := 12345

	// Get team members
	fmt.Printf("=== Team Members (ID: %d) ===\n", teamID)
	members, err := client.Teams.Team(teamID).Members(ctx)
	if err != nil {
		log.Printf("Failed to get team members: %v\n", err)
	} else {
		fmt.Printf("Team has %d members:\n", len(members.Data))
		for _, member := range members.Data {
			fmt.Printf("- %s (Rank Title: %s, Rank: %d)\n",
				member.Name, member.RankText, member.Points)
		}
	}

	// Get team activity for the last 30 days
	fmt.Printf("\n=== Team Activity ===\n")
	activity, err := client.Teams.Team(teamID).Activity(ctx, 30)
	if err != nil {
		log.Printf("Failed to get team activity: %v\n", err)
	} else {
		fmt.Printf("Recent team activities: %d entries\n", len(activity.Data))
		for i, act := range activity.Data {
			if i >= 3 { // Show only first 3 activities
				break
			}
			fmt.Printf("- %s: %s by %s\n", act.Date, act.Name, act.User.Name)
		}
	}

	// team info
	fmt.Printf("\n=== Team Info ===\n")
	info, err := client.Teams.TeamInfo(ctx, teamID)
	if err != nil {
		log.Printf("failed to get team info: %v\n", err)
	}else {
		fmt.Printf("Team Motto: %s", info.Data.Motto)
	}

	// team stats
	fmt.Printf("\n=== Team Stats ===\n")
	stats, err := client.Teams.Stats(ctx, teamID)
	if err != nil {
		log.Printf("failed to get team stats: %v\n", err)
	}else {
		fmt.Printf("team stats:\n Bloods: %d\n System Owns: %d\n User Owns: %d", stats.Data.FirstBloods, stats.Data.SystemOwns, stats.Data.UserOwns)
	}

}
