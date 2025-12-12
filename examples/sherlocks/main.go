package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gubarz/gohtb"
)

func main() {
	fmt.Println("=== Sherlocks example ===")
	ctx := context.Background()

	client, err := gohtb.New(os.Getenv("HTB_TOKEN"))
	if err != nil {
		log.Fatalln("Failed to create HTB client:", err)
	}
	fmt.Println("=== Getting all Sherlocks ===")
	sherlocks, err := client.Sherlocks.List().Results(ctx)
	if err != nil {
		log.Panicf("cannot fetch all sherlocks %v!", err)
	}
	for _, sherlock := range sherlocks.Data {
		fmt.Printf("Name: %s, State: %s\n", sherlock.Name, sherlock.State)
	}
}