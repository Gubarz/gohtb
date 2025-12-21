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

	// get Brutus Sherlock as example
	sherlock_id := 631
	sherlock := client.Sherlocks.Sherlock(sherlock_id)
	info, err := sherlock.Info(ctx) // grab info like avatar or name (basic info)
	if err != nil {
		log.Fatalln("Failed to get brutus sherlock:", err)
	}
	progress, err := sherlock.Progress(ctx) // sherlock progress
	if err != nil {
		log.Fatalln("Failed to get sherlock progress:", err)
	}
	playInfo, err := sherlock.Play(ctx) // play information like scenario or other info for playing
	if err != nil {
		log.Fatalln("Failed to get play info:", err)
	}
	fmt.Printf("Name: %s\n Percentage Completed: %.0f%%\n Description: %s\n", info.Data.Name, float64(progress.Data.Progress), playInfo.Data.Scenario)
}
