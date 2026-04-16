package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gubarz/gohtb"
)

func main() {
	ctx := context.Background()

	client, err := gohtb.New(os.Getenv("HTB_TOKEN"))
	if err != nil {
		log.Fatalln("Failed to create HTB client:", err)
	}

	// List active hard web challenges
	fmt.Println("=== Active Hard Challenges ===")
	challenges, err := client.Challenges.List().
		ByState("active").
		ByDifficulty("hard").
		SortedBy("rating").Descending().
		Results(ctx)
	if err != nil {
		log.Printf("Failed to get challenges: %v\n", err)
	} else {
		fmt.Printf("Found %d hard challenges:\n", len(challenges.Data))
		for i, challenge := range challenges.Data {
			if i >= 5 { // Show only first 5
				break
			}
			fmt.Printf("- %s (ID: %d, Category: %s, Rating: %.1f)\n",
				challenge.Name, challenge.Id, challenge.CategoryName, challenge.Rating)
		}

		if len(challenges.Data) > 0 {
			// Get detailed info for the first challenge
			firstChallenge := challenges.Data[0]
			fmt.Printf("\n=== Challenge Details: %s ===\n", firstChallenge.Name)

			info, err := client.Challenges.Challenge(firstChallenge.Id).Info(ctx)
			if err != nil {
				log.Printf("Failed to get challenge info: %v\n", err)
			} else {
				fmt.Printf("Description: %s\n", info.Data.Description)
				fmt.Printf("Points: %d\n", info.Data.Points)
				fmt.Printf("Difficulty: %s\n", info.Data.Difficulty)
				fmt.Printf("Solves: %d\n", info.Data.Solves)

				// Check if challenge has downloadable files
				if info.Data.FileName != "" {
					fmt.Printf("\n=== Downloading Challenge Files ===\n")
					resp, err := client.Challenges.Challenge(firstChallenge.Id).DownloadLink(ctx)
					if err != nil {
						log.Printf("Failed to get download link: %v\n", err)
					} else {
						fmt.Printf("Download URL: %s\n", resp.Data.Url)

						// Download the file
						downloadResp, err := http.Get(resp.Data.Url)
						if err != nil {
							log.Printf("Failed to download file: %v\n", err)
							return
						}

						downloadData, err := io.ReadAll(downloadResp.Body)
						if err != nil {
							log.Printf("Failed to read downloaded data: %v\n", err)
							return
						}
						downloadResp.Body.Close()

						// Save the file locally
						err = os.WriteFile(info.Data.FileName, downloadData, 0644)
						if err != nil {
							log.Printf("Failed to save file: %v\n", err)
						} else {
							fmt.Printf("Downloaded: %s (%d bytes)\n", info.Data.FileName, len(downloadData))
						}
					}
				} else {
					fmt.Println("No downloadable files for this challenge")
				}

				// Start challenge instance (if supported)
				fmt.Printf("\n=== Starting Challenge Instance ===\n")
				startResult, err := client.Challenges.Challenge(firstChallenge.Id).Start(ctx)
				if err != nil {
					log.Printf("Failed to start challenge: %v\n", err)
				} else {
					fmt.Printf("Start result: %s\n", startResult.Data.Message)
				}
			}
		}
	}
}
