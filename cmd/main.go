package main

import (
	"context"
	"fmt"
	"os"
	"github.com/dkaman/sumo-api-golang/pkg/sumoapi"
)

func main() {
	ctx := context.Background()
	client := sumoapi.NewClient(nil)
	rikishi, _, err := client.Rikishi.List(ctx)
	if err != nil {
		fmt.Printf("error with rikishi list %s\n", err)
		os.Exit(1)
	}

	for _, r := range rikishi {
		fmt.Println("rikishi: %s\n", r)

	}
}
