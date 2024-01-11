package main

import (
	"context"
	"fmt"
	"github.com/dkaman/sumo-api-golang/pkg/sumoapi"
	"os"
)

func main() {
	ctx := context.Background()
	client := sumoapi.NewClient(nil)
	// rikishi, _, err := client.Rikishi.List(ctx)
	// if err != nil {
	// 	fmt.Printf("error with rikishi list %s\n", err)
	// 	os.Exit(1)
	// }

	// for _, r := range rikishi {
	// 	fmt.Println("rikishi: %s\n", r)

	// }

	// r, _,err := client.Rikishi.Get(ctx, 8890)
	// if err != nil {
	// 	fmt.Printf("%s\n", err)
	// 	os.Exit(1)

	// }
	// fmt.Printf("r = %s\n", r)

	stats, _, err := client.Rikishi.Stats(ctx, 12)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(stats)

	matches, _, err := client.Rikishi.Matches(ctx, 12)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, m := range matches {
		fmt.Println(m)
	}

	matchup, _, err := client.Rikishi.Matchup(ctx, 12, 72)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(matchup)
}
