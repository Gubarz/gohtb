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

}
