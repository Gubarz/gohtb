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

	// Get Machine Info
	fmt.Println("=== Machine Info ===")
	info, err := client.Machines.Machine(660).Info(ctx)
	if err != nil {
		log.Printf("Failed to get machine info: %v\n", err)
	} else {
		fmt.Printf("Name: %s (ID: %d, OS: %s, Difficulty: %s)\n",
			info.Data.Name, info.Data.Id, info.Data.Os, info.Data.DifficultyText)

	}
	// getting all active and unreleased
	fmt.Println("=== Active and Unreleased Machines Info ===")
	infos, err := client.Machines.ListUnreleased().AllResults(ctx)
	if err != nil {
		log.Panicf("Failed to get unreleased machines: %v\n", err)
	} else {
		for _, machine := range infos.Data {
			fmt.Printf("Name: %s (OS: %s Difficulty: %s)\n", machine.Name, machine.Os, machine.DifficultyText)
		}
	}
	fmt.Println("=== Now Active Machines ===")
	infos, err = client.Machines.ListActive().AllResults(ctx)

	if err != nil {
		log.Panicf("Failed to get active machines: %v\n", err)
	} else {
		for _, machine := range infos.Data {
			fmt.Printf("Name: %s (OS: %s Difficulty: %s)\n", machine.Name, machine.Os, machine.DifficultyText)
		}
	}
	fmt.Println("=== Testing retired machines ===")
	infos, err = client.Machines.ListRetired().Results(ctx)
	if err != nil {
		log.Panicf("Failed to get retired machines: %v\n", err)
	} else {
		for _, machine := range infos.Data[:10] {
			fmt.Printf("Name: %s (OS: %s Difficulty: %s State: %s)\n", machine.Name, machine.Os, machine.DifficultyText, machine.State)
		}
	}


	// Get unreleased machine from all Results without any filtering
	fmt.Println("=== Getting unreleased machines using new `List().AllResults(ctx)` ===")
	v5_search, err := client.Machines.List().AllResults(ctx)
	if err != nil {
		log.Printf("Failed to get machine info: %v\n", err)
	} else {
		for _, info := range v5_search.Data {
			if info.State == "unreleased" {
				fmt.Printf("Name: %s (ID: %d, OS: %s, Difficulty: %s State: %s)\n",
					info.Name, info.Id, info.Os, info.DifficultyText, info.State)
			}
		}

	}

}
