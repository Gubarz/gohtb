package main
import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gubarz/gohtb"
)


func main(){
	ctx := context.Background()

	client, err := gohtb.New(os.Getenv("HTB_TOKEN"))
	if err != nil {
		log.Fatalln("Failed to create HTB client:", err)
	}
	ranks, err := client.Rankings.Rankings().Users(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(ranks.Data))
}